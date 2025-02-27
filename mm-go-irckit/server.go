package irckit

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/42wim/matterircd/bridge"
	"github.com/sorcix/irc"
)

var ErrHandshakeFailed = errors.New("handshake failed")

var defaultVersion = "go-irckit"

const handshakeMsgTolerance = 20

// ID will normalize a name to be used as a unique identifier for comparison.
func ID(s string) string {
	return strings.ToLower(s)
}

type Prefixer interface {
	// Prefix returns a prefix configuration for the origin of the message.
	Prefix() *irc.Prefix
}

type Server interface {
	Prefixer

	// Name of the server (usually hostname).
	Name() string

	// Motd is the Message of the Day for the server.
	Motd() []string

	// Connect starts the handshake for a new user, blocks until it's completed or failed with an error.
	Connect(*User) error

	// Quit removes the user from all the channels and disconnects.
	Quit(*User, string)

	// HasUser returns an existing User with a given Nick.
	HasUser(string) (*User, bool)

	// HasUserID returns an existing User with a given ID
	HasUserID(string) (*User, bool)

	// RenameUser changes the Nick of a User if the new name is available.
	// Returns whether the rename was was successful.
	RenameUser(*User, string) bool

	// Channel gets or creates a new channel with the given name and Id.
	Channel(string) Channel

	// HasChannel returns an existing Channel with a given name.
	HasChannel(string) (Channel, bool)

	// UnlinkChannel removes the channel from the server's storage if it
	// exists. Once removed, the server is free to create a fresh channel with
	// the same ID. The server is not responsible for evicting members of an
	// unlinked channel.
	UnlinkChannel(Channel)

	Add(u *User) bool
	BatchAdd(users []*User)
	Handle(u *User)
	Logout(u *User)
	ChannelCount() int
	UserCount() int
	EncodeMessage(u *User, cmd string, params []string, trailing string) error
}

// ServerConfig produces a Server setup with configuration options.
type ServerConfig struct {
	// Name is used as the prefix for the server.
	Name string
	// Version string of the server (default: go-irckit).
	Version string
	// Motd is the message of the day for the server, list of message lines where each line should be max 80 chars.
	Motd []string
	// InviteOnly prevents regular users from joining and making new channels.
	InviteOnly bool
	// MaxNickLen is the maximum length for a NICK value (default: 32)
	MaxNickLen int

	// DiscardEmpty setting will start a goroutine to discard empty channels.
	DiscardEmpty bool
	// NewChannel overrides the constructor for a new Channel in a given Server and Name.
	NewChannel func(s Server, channelId string, name string, service string, modes map[string]bool) Channel
	// Commands is the handler registry to use (default: DefaultCommands())
	Commands Commands
}

func (c ServerConfig) Server() Server {
	if c.NewChannel == nil {
		c.NewChannel = NewChannel
	}
	if c.Commands == nil {
		c.Commands = DefaultCommands()
	}
	if c.Version == "" {
		c.Version = defaultVersion
	}
	if c.Name == "" {
		c.Name = "go-irckit"
	}
	if c.MaxNickLen == 0 {
		c.MaxNickLen = 32
	}

	srv := &server{
		config:   c,
		users:    map[string]*User{},
		channels: map[string]Channel{},
		created:  time.Now(),
		commands: c.Commands,
	}

	return srv
}

// NewServer creates a server.
func NewServer(name string) Server {
	return ServerConfig{Name: name}.Server()
}

type server struct {
	config   ServerConfig
	created  time.Time
	commands Commands

	u *User
	sync.RWMutex
	users    map[string]*User
	channels map[string]Channel
}

func (s *server) Name() string {
	return s.config.Name
}

func (s *server) Motd() []string {
	return s.config.Motd
}

func (s *server) Close() error {
	// TODO: Send notice or something?
	// TODO: Clear channels?
	s.Lock()
	for _, u := range s.users {
		u.Close()
	}
	s.Unlock()
	return nil
}

// Prefix returns the server's command prefix string.
func (s *server) Prefix() *irc.Prefix {
	return &irc.Prefix{Name: s.config.Name}
}

// HasUser returns whether a given user is in the server.
func (s *server) HasUser(nick string) (*User, bool) {
	s.RLock()
	defer s.RUnlock()
	for _, u := range s.users {
		u := u
		if u.Nick == nick {
			return u, true
		}
	}
	//	u, exists := s.users[ID(nick)]
	return nil, false
}

func (s *server) HasUserID(userID string) (*User, bool) {
	s.RLock()
	u, exists := s.users[strings.ToLower(userID)]
	s.RUnlock()
	return u, exists
}

// Rename will attempt to rename the given user's Nick if it's available.
func (s *server) RenameUser(u *User, newNick string) bool {
	if len(newNick) > s.config.MaxNickLen {
		newNick = newNick[:s.config.MaxNickLen]
	}

	s.Lock()
	if _, exists := s.users[ID(newNick)]; exists {
		s.Unlock()
		s.EncodeMessage(u, irc.ERR_NICKNAMEINUSE, []string{newNick}, "Nickname is already in use")
		return false
	}

	delete(s.users, u.ID())
	oldPrefix := u.Prefix()
	u.Nick = newNick
	s.users[u.ID()] = u
	s.Unlock()

	changeMsg := &irc.Message{
		Prefix:  oldPrefix,
		Command: irc.NICK,
		Params:  []string{newNick},
	}
	u.Encode(changeMsg)
	for _, other := range u.VisibleTo() {
		other.Encode(changeMsg)
	}
	return true
}

// HasChannel returns whether a given channel already exists.
func (s *server) HasChannel(channelID string) (Channel, bool) {
	s.RLock()
	ch, exists := s.channels[channelID]
	s.RUnlock()
	return ch, exists
}

// Channel returns an existing or new channel with the give name.
func (s *server) Channel(channelID string) Channel {
	s.Lock()
	ch, ok := s.channels[channelID]
	if !ok {
		service := s.u.br.Protocol()
		name := s.u.br.GetChannelName(channelID)

		info, err := s.u.br.GetChannel(channelID)
		if err != nil {
			// don't error on our special channels
			if !strings.HasPrefix(channelID, "&") {
				logger.Errorf("didn't find channel %s: %s", channelID, err)
			}
			info = &bridge.ChannelInfo{}
		}

		modes := make(map[string]bool)
		modes["p"] = info.Private

		newFn := s.config.NewChannel
		ch = newFn(s, channelID, name, service, modes)

		logger.Debugf("new channel id: %s, name: %s", channelID, name)

		s.channels[channelID] = ch
		s.channels[name] = ch
		s.Unlock()
	} else {
		s.Unlock()
	}
	return ch
}

// UnlinkChannel unlinks the channel from the server's storage, returns whether it existed.
func (s *server) UnlinkChannel(ch Channel) {
	s.Lock()
	chStored := s.channels[ch.String()]
	r := chStored == ch
	if r {
		delete(s.channels, ch.String())
		delete(s.channels, ch.ID())
	}
	s.Unlock()
}

// Connect starts the handshake for a new User and returns when complete or failed.
func (s *server) Connect(u *User) error {
	err := s.handshake(u)
	if err != nil {
		u.Close()
		return err
	}
	go s.handle(u)
	return nil
}

// Quit will remove the user from all channels and disconnect.
func (s *server) Quit(u *User, message string) {
	go u.Close()
	s.Lock()
	delete(s.users, u.ID())
	s.Unlock()

	u.br.Logout()
}

// Len returns the number of users connected to the server.
func (s *server) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.users)
}

func (s *server) welcome(u *User) error {
	err := u.Encode(
		&irc.Message{
			Prefix:   s.Prefix(),
			Command:  irc.RPL_WELCOME,
			Params:   []string{u.Nick},
			Trailing: fmt.Sprintf("Welcome! %s", u.Prefix()),
		},
		&irc.Message{
			Prefix:   s.Prefix(),
			Command:  irc.RPL_YOURHOST,
			Params:   []string{u.Nick},
			Trailing: fmt.Sprintf("Your host is %s, running version %s", s.config.Name, s.config.Version),
		},
		&irc.Message{
			Prefix:   s.Prefix(),
			Command:  irc.RPL_CREATED,
			Params:   []string{u.Nick},
			Trailing: fmt.Sprintf("This server was created %s", s.created.Format(time.UnixDate)),
		},
		&irc.Message{
			Prefix:   s.Prefix(),
			Command:  irc.RPL_MYINFO,
			Params:   []string{u.Nick},
			Trailing: fmt.Sprintf("%s %s o o debugmode %t", s.config.Name, s.config.Version, IsDebugLevel()),
		},
		&irc.Message{
			Prefix:   s.Prefix(),
			Command:  irc.RPL_LUSERCLIENT,
			Params:   []string{u.Nick},
			Trailing: fmt.Sprintf("There are %d users and 0 services on 1 servers", s.Len()),
		},
	)
	if err != nil {
		return err
	}
	// Always include motd, even if it's empty? Seems some clients expect it (libpurple?).
	return CmdMotd(s, u, nil)
}

func (s *server) EncodeMessage(u *User, cmd string, params []string, trailing string) error {
	return u.Encode(&irc.Message{
		Prefix:   s.Prefix(),
		Command:  cmd,
		Params:   params,
		Trailing: trailing,
	})
}

func (s *server) Handle(u *User) {
	s.handle(u)
}

func (s *server) handle(u *User) {
	var partMsg string
	defer s.Quit(u, partMsg)
	for msg := range u.DecodeCh {
		if msg == nil {
			// Ignore empty messages
			continue
		}
		go func(msg *irc.Message) {
			err := s.commands.Run(s, u, msg)
			logger.Debugf("Executed %#v %#v", msg, err)
			if err == ErrUnknownCommand {
				// TODO: Emit event?
			} else if err != nil {
				logger.Errorf("handler error for %s: %s", u.ID(), err.Error())
			}
		}(msg)
	}
}

func (s *server) BatchAdd(users []*User) {
	s.Lock()
	defer s.Unlock()

	for _, u := range users {
		id := u.ID()
		if _, exists := s.users[id]; !exists {
			s.users[id] = u
		}
	}
}

func (s *server) Add(u *User) (ok bool) {
	return s.add(u)
}

func (s *server) add(u *User) (ok bool) {
	s.Lock()
	defer s.Unlock()

	id := u.ID()
	if _, exists := s.users[id]; exists {
		return false
	}

	s.users[id] = u
	return true
}

//nolint:goconst
func (s *server) handshake(u *User) error {
	// Assign host
	u.Host = u.ResolveHost()
	go u.Decode()

	// Consume N messages then give up.
	i := handshakeMsgTolerance
	// Read messages until we filled in USER details.
outerloop:
	for {
		select {
		case msg := <-u.DecodeCh:
			// fmt.Printf("in handshake %#v\n", msg)
			i--
			// Consume N messages then give up.
			if i == 0 {
				break outerloop
			}
			if msg == nil {
				// Empty message, ignore.
				continue
			}

			// apparently NICK message can have a : prefix on connection
			// https://github.com/42wim/matterircd/issues/32
			if (msg.Command == irc.NICK || msg.Command == irc.PASS) && msg.Trailing != "" {
				msg.Params = append(msg.Params, msg.Trailing)
			}
			if len(msg.Params) < 1 {
				continue
			}

			switch msg.Command {
			case irc.NICK:
				u.Nick = msg.Params[0]
			case irc.USER:
				u.User = msg.Params[0]
				u.Real = msg.Trailing
			case irc.PASS:
				u.Pass = msg.Params
			case irc.JOIN:
				s.EncodeMessage(u, irc.ERR_NOTREGISTERED, []string{"*"}, "Please register first")
			}

			if u.Nick == "" || u.User == "" {
				// Wait for both to be set before proceeding
				continue
			}
			if len(u.Nick) > s.config.MaxNickLen {
				u.Nick = u.Nick[:s.config.MaxNickLen]
			}

			ok := s.add(u)
			if !ok {
				s.EncodeMessage(u, irc.ERR_NICKNAMEINUSE, []string{u.Nick}, "Nickname is already in use")
				continue
			}
			s.u = u

			err := s.welcome(u)
			if err == nil && u.Pass != nil {
				service := "mattermost"
				if len(u.Pass) == 1 {
					service = "slack"
				}
				login(u, &User{
					UserInfo: &bridge.UserInfo{
						Nick: service,
						User: service,
						Real: service,
						Host: "service",
					},
					channels: map[Channel]struct{}{},
				},
					u.Pass,
					service)
			}

			return err
		case <-time.After(10 * time.Second):
			return ErrHandshakeFailed
		}
	}
	return ErrHandshakeFailed
}

func (s *server) Logout(user *User) {
	channels := user.Channels()
	for _, ch := range channels {
		for _, other := range ch.Users() {
			s.Lock()
			delete(s.users, other.ID())
			s.Unlock()
		}
		ch.Part(user, "")
		ch.Unlink()
	}
}

func (s *server) ChannelCount() int {
	return len(s.channels)
}

func (s *server) UserCount() int {
	return len(s.users)
}
