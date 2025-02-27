# v0.24.1

## Enhancement

- mattermost: Don't show warning if lastViewedAt state file doesn't exist, also create it early (#423)
- slack: update vendored library

## Bugfix

- mattermost: Fix ordering of events (such as code blocks) (#419)
- mattermost: Save last viewed at state on logout and only use stored last viewed at if newer than what the server knows (#422)
- slack: Make sure receiving item message isn't nil (slack). Fixes #428 (#429)

This release couldn't exist without the following contributors:
@hloeung

# v0.24.0

## New features

- mattermost: Add context for reactions (#408)
- mattermost: Add option to hide reactions; also hide our own (#415)
- mattermost: Fix channel mentions and add option to disable channel mentions as IRC NOTICEs (mattermost) (#418)
- mattermost: Show Mattermost message/thread IDs for file attachments (#396) (was already in v0.23.1)
- mattermost: Show scrollback in same channel/query window (mattermost) (#393) (was already in v0.23.1)

## Enhancements

- mattermost: Save lastViewedAt on topic change events (#417)
- slack: Save message IDs and use them for local echo prevention (#409)

## Bugfix

- general: Set correct IRC PRIVMSG nick when receiving direct messages (#407)
- general: Fix buffer flushing yet again (#411)
- mattermost: Handle SIGSEGV nil pointer for event.ParentUser (#412)

This release couldn't exist without the following contributors:
@hloeung, @the-real-ed, @Aketzu

# v0.23.1

## Bugfix

- slack: Use conversation API for DM (slack). (#402) Fixes #400
- slack: Fix regression, increase user pag to 1000 (slack) #403
- general: Fixed buffer flushing (#398)

This release couldn't exist without the following contributors:
@hloeung

# v0.23.0

Big shout-out to @hloeung for doing most of the work for this release!

## New features

- mattermost: Replay only messages not seen; also make it clear when replaying #373
- mattermost: new ForceAntiIdle option - Allow forcing anti-idle #369
- mattermost: Allow adding or removing reactions #386
- mattermost: Add scrollback support for direct messages (DMs). Fixes #385 #388

## Enhancements

- mattermost: Add prefix showing messages are replies to threads when using Mattermost ThreadContext
- mattermost: Save lastViewedAt and load on start up (#313) #368
- mattermost: Add logging when replaying logs #377
- general: Handle events after initialization use buffered channel #383

## Bugfix

- general: Flush buffer on ctcp action. Fixes #367 #371
- general: Sanitize nicks, replace invalid chars with dash #374
- general: Sanitize nick before sending PRIVMSG irc command #378
- general: Flush buffers on reactions, allows adding various reactions in sequence, also replies to threads or message modifications #387
- general: Fix fd leak - Timeout unused handshake after 10 seconds. Fixes #390 #392
- general: Space separate 'download file' #394
- mattermost: Fix panic with no postlist returned. Fixes #380 #384
- mattermost: Fix permissions issues when replying to threads in DM #389

This release couldn't exist without the following contributors:
@hloeung, @dfskoll

# v0.22.0

## New features

- mattermost: Allow modifying or deleting last message and by Mattermost message ID (mattermost) (#353)
- mattermost: Allow replying to your own last message to start new threads (mattermost) (#362)
- mattermost: Support 2fa authentication (mattermost) (#355)

## Enhancements

- general: Log useful timestamp (#364)
- mattermost: Save last viewed at state and use for replaying logs (#361)
- mattermost: Add optional Unicode ellipsis to save on screen real estate (#358)
- mattermost: Allow showing replies with Prefix/SuffixContext (#351)

## Bugfix

- mattermost: Do not show replies when hidereplies is enabled (mattermost) (#352)

This release couldn't exist without the following contributors:
@hloeung, @eyenx

# v0.21.0

## New features

- general: Reload certs using SIGHUP. Fixes #110
- general: Add a tlskey and tlscert configuration option
- mattermost: Support replying to old mattermost threads (#344)
- mattermost: Allow showing parent thread or message ID with PrefixContext / SuffixContext (#348)

## Bugfix

- general: Disable watchconfig on illumos. Fixes #345
- mattermost: Use override_username when available in backlog (mattermost) (#343)

This release couldn't exist without the following contributor:
@hloeung

# v0.20.2

## New features

- mattermost: Add ShortenReplies (#322). see matterircd.toml.example

## Bugfix

- general: Ignore private messages when not logged in (#331)
- general: Lock updateCounter against concurrent updates (#328)
- general: Ignore empty channels and non-channels for TOPIC (#339)
- mattermost: Move antiIdle to matterclient package. Fixes #327 (#329)
- mattermost: Update user after error checking (#330)
- mattermost: Fix nil User when using tokens (mattermost) (#334)
- mattermost: Use channel id from CreateDirectChannel (mattermost). Fixes #336 (#337)
- slack: Create a fake user when getting a botmessage (slack) (#333)
- slack: Fix missing botname, fix attachments without text (slack) (#335)

This release couldn't exist without the following contributors:
@Aketzu, @jetpackdanger

# v0.20.1

## New features

- slack: support SSO/SAML logins using xoxc-token and cookies. See <https://github.com/42wim/matterircd#slack-sso-login--xoxc-tokens>
- mattermost: Support for adding message/thread context at the end (#319). see matterircd.toml.example

## Bugfix

- general: Fix panic on nick change when not logged in
- general: Fix panic on sending to messages/users channels
- mattermost: Do not check alive while reconnecting. Fixes #318
- mattermost: Use always other user as channelID in DM for prefixcontext. Fixes #317
- slack: fix possible panic when self-message #321

This release couldn't exist without the following contributors:
@hloeung, @42wim

# v0.20.0

The refactor edition.
Be sure to look into the new [prefixcontext.md](https://github.com/42wim/matterircd/blob/master/prefixcontext.md) option for mattermost.

Thanks a lot to @cjwatson, @muesli, @hloeung, @vmiklos and @guilhermepiccoli for contributing and testing this release!

## Breaking changes

### Commandline switches

- Switched to viper for cmdline parsing, which does not support "short" flags. You'll need to use `--flag` instead of `-flag`. Eg `./matterircd --debug`
- Bridge specific configuration is now only in configuration file. This means the following flags have been removed: `-restrict`,`-mmteam`,`-mmserver`,`-mminsecure`,`-mmskiptlsverify`. You can set those in `matterircd.toml`, see the example file.

### Config changes

- `BlacklistUser` feature for slack has been renamed to `DenyUsers`.
- `JoinMpImOnTalk` feature has been renamed to `JoinDM` and is available for slack/mattermost
- `JoinInclude`, `JoinExclude` now support regexp (see matterircd.toml.example)

## New features

- general: New option `trace` to give even more output as `debug`. See matterircd.toml.example.
- general: New option `gops` to activate [gops](https://github.com/google/gops) for debugging.
- general: Allow binding to a Unix socket #276.
- mattermost: Add prefixcontext option which supports edits/deletes/reactions/threads (see matterircd.toml.example) and ([prefixcontext.md](https://github.com/42wim/matterircd/blob/master/prefixcontext.md))
- mattermost: Add option to use Nickname instead of Username #273 (See matterircd.toml.example).
- mattermost: Add option to disable showing replies/parent posts #283 (See matterircd.toml.example).
- mattermost: `JoinDM` option to enable groups/dm joining on startup.
- mattermost: Add thread replies to prefixcontext.
- mattermost: Add support for thread reply using prefixcontext.
- mattermost: Add support for deleting/modifying own messages.
- mattermost: Add support for "words that trigger mentions" (See matterircd.toml.example)
- mattermost/slack: `JoinOnly` option to only join these specific channels (see matterircd.toml.example).
- mattermost/slack: `JoinInclude`, `JoinExclude` now support regexp (see matterircd.toml.example).

## Enhancement

- general: Refactor using interfaces, will make it easier to update and add new bridges.
- general: massive speedups on joining large channels.
- general: less memory/CPU usage on large slacks/mattermost.
- general: Break longer messages at word boundaries #270.
- general: Move matterclient to pkg and refactor.
- general: Batch lastviewed updates.
- mattermost: Handle ratelimiting on frequent used paths.
- mattermost: Rewritten mattermost connection logic fixing goroutine leaks.
- mattermost: make `@ALL` messages also notices #288.
- mattermost: add support for updateuser event (realtime nick changes).
- mattermost: private channels are now shown with +p in irc
- mattermost: Convert /me commands from mattermost to irc /me. Closes #281
- mattermost: Handle deleted messages in mattermost.
- mattermost: Add reaction events (using prefixContext).
- slack: speed-up on large slack installations.

## Bugfix

- general: Refactor DMs and fix messages to self. Fixes #293
- general: Fix topic messages and check permissions. Closing #294
- general: Fix too many replay messages.
- mattermost: Changing topic also changes channel display name #284.
- mattermost: Images/links in private messages now are on the correct channel.
- mattermost: Ignore user join messages #280
- mattermost: Make mattermost away work. Fix #277
- mattermost: Fix goroutine leaks on multiple login/logout
- mattermost: Fix crash on quit without login #300
- slack: fix double messages on irc because of slack API changes, now using slack blocks.

# v0.19.4

## Bugfix

- slack: fix regression with slack library (#264)
- slack: fix an unexpected panic (#263)

# v0.19.3

## Enhancement

- general: Add UPDATELASTVIEWED command, and make DisableAutoView work consistently (#255)
- slack: Handle message edits and deletion (#260)
- slack: Add handling of reactions, stars and pins (#229)

## Bugfix

- mattermost: Fix a panic #247
- mattermost: Fixes incorrect users because of paging. #244
- mattermost: Fix outdated channel issue
- mattermost: Add paging so we can see > 200 users in a channel #248
- mattermost: Fix expired session panic #259
- general: Fix datarace #246
- general: Fix empty JoinInclude
- general: Fix panic #257

This release couldn't exist without the following contributors:
@Aketzu, @bucko909, @42wim

# v0.19.2

## Enhancement

- general: Add a default value matterirc.toml for the '-conf' flag (#240)
- slack: library updated
- mattermost: library updated
- mattermost: Add support for channel created/deleted events

## Bugfix

- mattermost: Remove ourselves from the channel when removed in mattermost. Fixes #233
- mattermost: Add/remove ourselves to the channel if we join using the GUI. #239
- mattermost: Update topics in mattermost. Closes #241
- mattermost: Fix pastes and attachments in direct message. Closes #228
- mattermost: Update channels if not known on join yet

# v0.19.1

## New features

- mattermost: Added support for disabling of automatic view flag updates (#226). See DisableAutoView in matterircd.toml.example
- slack: Add message showing enhancements and add slackbot to all channels (#230)

## Bugfix

- general: Fix tight loop (100% CPU). Closes #231

# v0.19.0

## New features

- irc: Add support for spoofing query messages. #195
  - You can now see your own messages you've typed on slack/mattermost web in irc
- irc: Add PasteBufferTimeout option (send ascii-art!)
  _ See matterircd.toml.example for an example.
  _ PasteBufferTimeout specifies the amount of time in milliseconds that messages get kept in matterircd internal buffer before being sent to
  mattermost or slack. Messages that will be received in this time will be concatenated together
  So this can be used to paste stuff like ascii-art or code.
  Default 0 (is disabled)
  Depending on how fast you type 2500 is a good number

## Bugfix

- slack: Correctly handle different nick and username #203
- slack: Ignore channel join messages #198

# v0.18.4

## Bugfix

- general: fix cli args not override configuration file #205
- mattermost: support multi DM-groups correctly #209
- mattermost: add correct support for personal tokens #208
  `Use /msg mattermost login <server> <team> <login> token=<yourtoken>`
- mattermost: Fix JoinInclude / JoinExclude logic when joining/parting channels. Also support #team/channel
- mattermost: Fix issue with empty channelname
- mattermost: Fix re-login on session expiry

# v0.18.3

## Bugfix

- slack: api changed, show uploaded files again

Because of changes in slack API and the forced use of pagination, big channels with lots of users can take a while to load.

# v0.18.2

## Bugfix

- slack: fix panic on websocket bug #189, #196

# v0.18.1

## New features

- mattermost: support mattermost 5.x

# v0.18.0

## New features

- general: Add debugmode true/false message in banner
- irc: Add PrefixMainTeam setting. Also set the main team name as prefix in the irc-channel. See matterircd.toml.example
- slack: Add support for login <team> <user> <pass> for slack (as addition to login <token>)

## Bugfix

- mattermost: update channels when adding/removing users to new channel. Alsso join channel when we are added. Closes 42wim/matterircd#178
- irc: fix NAMES reply to send entire member list
- irc: Use service for on-join topic changes (instead of your own username)
- irc: Handle \r, ACTION and colour sanitization everywhere
- irc: Fix concurrent map read/write. Closes 42wim/matterircd#188
- slack: Make sure to return for not implemented functions in slack. Closes 42wim/matterircd#186
- slack: Replace spaces to underscore in botnames. Closes 42wim/matterircd#184

# v0.17.3

## Bugfix

- slack: Fix issues with bots with spaces in the name
- mattermost: Actually join/remove users to channel when they join, not when they say something #113
- mattermost: Join/remove users when they're added by someone else. Use a system message to show this #175

# v0.17.2

## Bugfix

- mattermost: Fix message looping issue

# v0.17.1

## New features

- general: enable login via irc PASS command during handshake instead of PRIVMSG

## Bugfix

- mattermost: Update GetFileLinks to API_V4
- slack/mattermost: Fix issue with matterircd users not being able to chat to eachother
- slack: Do not join channels for single direct messages (slack)
- slack: Split fallback messages on newline (slack)

# v0.17.0

## New features

- general: mattermost configuration settings need to be migrated to `[mattermost]` settings. See matterircd.toml.example
- slack: Add BlackListUser config setting. Blacklist users from connecting to slack. See matterircd.toml.example
- slack: Add JoinMpImOnTalk config setting. Only join MultiPerson IM when someone talks. See mattericd.toml.example
- slack: Add Restrict config setting. Only allowed to connect to specified slack teams. See matterircd.toml.example
- slack: Add UseDisplayName config setting. If displayname is set, the message will be prepended with `<displayname>`. See matterircd.toml.example
- slack: also show messages from bots
- slack: reconnect on disconnects

## Bugfix

- &users join speedup on teams with massive amount of users (tested on 26k users)
- Only allow 1 login/logout in progress
- slack: Fix on-join race condition
- slack: Always add yourself to your channels (fixes problem with > 500 users channels)
- slack: remove carriage returns from topic
- slack: Autojoin new channel/group when invited
- slack: Join MpIm channel if we haven't joined it yet

# v0.16.8

## Bugfix

- Fix newlines in topic to client. Closes 42wim/matterircd#163
- Remove double part messages. Closes 42wim/matterircd#156
- mattermost: Fix away (mattermost). Closes 42wim/matterircd#165
- slack: Look into attachments if message is empty (slack). Fixes 42wim/matterircd#160

# v0.16.7

## Bugfix

- mattermost: Update lastview every 60 second as antiIdle (mattermost). Closes 42wim/matterircd#147
- slack: Make sure we only can execute comands if login is fully done. Closes 42wim/matterircd#154
- slack: Add invite support (slack). Closes 42wim/matterircd#159
- slack: Add list support (slack)
- slack: Fix channel parts for private channels
- slack: Add kick support (slack)
- slack: Allow to mention yourself.. Closes 42wim/matterircd#157

# v0.16.6

## Bugfix

- Strip IRC colors. Closes 42wim/matterircd#149
- slack: Replace mentions, channels, etc.. #13
- slack: Use displayname if possible (instead of name)
- slack: Enable LinkNames. (@here,@all) Closes 42wim/matterircd#152
- slack: Add topic change support (slack). Closes 42wim/matterircd#151
- mattermost: Print list of valid team names when team not found (42wim/matterbridge#390)

# v0.16.5

## New features

- Add support for private channels in slack #142

## Bugfix

- Slack: fixes join/parts #143, #146
- Slack: fixes away #144

# v0.16.4

## Bugfix

- Fix some messages going to &messages #140

# v0.16.3

## Bugfix

- Fix crash on /nick change when not logged in #141

# v0.16.2

## Bugfix

- Remove crash on channel lookup of private messages

# v0.16.1

## Bugfix

- Remove debug code which could cause a crash
- Only append channel name to sender once in &messages

# v0.16.0

## New features

- `-conf` option (for a config file). See https://github.com/42wim/matterircd/blob/master/matterircd.toml.example for an example. Thanks @slowbro for this PR.
- New config file options

  - JoinExclude: an array of channels that won't be joined on IRC.
    Messages that get sent to unjoined channels (but you're joined on mattermost) will
    get sent to the &messages channel.
    You can still /JOIN exclude channels.

  JoinExclude = ["#town-square","#boringchannel"]

  - JoinInclude: an array of channels that only will be joined on IRC.
    If it's empty, it means all channels get joined (except those defined in JoinExclude)
    Messages that get sent to unjoined channels (but you're joined on mattermost) will
    get sent to the &messages channel.

  JoinInclude = ["#devops"]

  - PartFake: a bool that defines if you do a /LEAVE or /PART on IRC it will also
    actually leave the channel on mattermost.
    Default false

  PartFake = true

- don't log passwords used with 'mattermost' and 'slack'. Closes #73

## Bugfix

- Already read messages are replayed again and again #130
- Update to latest mattermost (4.6) libs
- Deprecated flags `-bindinterface` and `-port` removed

# v0.15.0

## New features

- Support mattermost 4.2 and higher (4.x) (use mattermost v4 API)
- Add -mmskiptlsverify option to skip TLS certificate checks on mattermost

## Enhancements

- Display nickname, if set #120
- Replace IRC parsing function with shellwords like function to allow for passwords with spaces. (#8)

# v0.14.0

## New features

- Support mattermost 4.1
- Add initial support for slack

# v0.13.0

## New features

- Support mattermost 4.0

## Enhancement

- Show slack attachments if they have a fallback/contain text
- Show links even when public links are disabled #105
- Edited messages now have "(edited)" appended
- `-bind ""` now disables non-tls port-binding when you have `-tlsbind` specified #109

## Bugfix

- Long messages from mattermost will be split in multiple smaller messages #103
- Fix join/leave messages for recent mattermost versions #113, #104
- Ignore messages sent to &users #108
- Ignore posts that have a reaction (emoji) added #111

# v0.12.0

(thanks to @recht matterircd fork)

## New features

- Add KICK support
- Add INVITE support
- Also relay edited messages

## Enhancement

- Show the original message/author after replied messages
- Print timestamp of replayed messages
- Faster startup (joining channels)

## Bugfix

- Do not clear topic on empty /TOPIC command
- Fix various possible panics

# v0.11.6

## New features

- Support mattermost 3.10.0

# v0.11.5

## Bugfix

- Fix crash #97

# v0.11.4

## New features

- Support mattermost 3.9.0

## Enhancement

- Use props instead of ZWSP to check if message comes from matterircd (no more spaces after messages from matterircd -> mattermost)

# v0.11.3

## New features

- Support mattermost 3.7.0 and 3.8.0

## Bugfix

- Make public links (pasted images/attachments) work again

## New features

- Add support for public links in SCROLLBACK and SEARCH

# v0.11.2

## Bugfix

- Use correct status for /WHO #81
- Fix ISON for misbehaving clients #78

## New features

- Support mattermost 3.6.0

# v0.11.1

## Bugfix

- Fix crash on new channel joins #77
- Fix crash on channel join/leaves of other mattermost users #77

# v0.11.0

## New features

- Support mattermost 3.5.0
