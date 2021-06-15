# IMS
Incident management tool with slack.

# (WIP)

# Solution

As the center of communication at work has been replaced from e-mail to chat, you may have changed the alert notification destination of your monitoring tool to chat.
But hasn't it changed that people are reading the messages and making decisions and responses?
I'm using chat for a change!
This tool enables easy incident management through chat, accelerating ChatOps!

# Feature

This tool has two major functions.
1. check the posted messages according to the rules, label them if they fit the rules, and repost them to the incidents management channel.

![2](https://user-images.githubusercontent.com/22161385/117542584-2f2aaf00-b054-11eb-8405-570a21570f9c.png)

2. Check if there are any reactions in the incidents management channel, and output the ones that are not there.

![1](https://user-images.githubusercontent.com/22161385/117542581-2d60eb80-b054-11eb-893e-d61355040935.png)

In other words, you can use the following

1. Invite this tool to the channel you are throwing the monitoring message into. The tool checks in the messages
2. Periodically, the incidents will run and display a list of messages that have not been reacted to, so check the unacted ones and leave a history of your responses in the thread.

```
 - - 2021/05/08 23:13:15 - -
[message] [Rule1] Error test2 [date] 2021/05/05 12:47:27
```

This makes it possible to

1. Identify unanswered alerts
2. Filter for known messages

Add the addressed messages to the configuration so that the cycle of improvement can continue.
All this can be done on slack!

# installation

If you want to put it under the path, you can use the following.

```
go get github.com/yasutakatou/IMS
```

If you want to create a binary and copy it yourself, use the following.

```
git clone https://github.com/yasutakatou/IMS
cd IMS
go build .
```

[or download binary from release page.](https://github.com/yasutakatou/IMS/releases)
save binary file, copy to entryed execute path directory.

# uninstall

delete that binary. del or rm command. (it's simple!)

# set up

Please follow the steps below to set up your environment.

1. set tool like bot. 
- goto [slack api](https://api.slack.com/apps)
- Create New(an) App
	- define (Name)
	- select (Workspce)
	- Create App
- App-Level Tokens
	-Generate Token and Scopes
	- define (Name)
	- Add Scope
		- connections:write
	- Generate
		- Make a note of the token that begins with xapp-.
	- Done
- Socket Mode
	- Enable Socket Mode
		- On
- OAuth & Permissions
	- Scopes
	- Bot Token Scopes
		- channels:history
		- chat:write
		- reactions:write
		- users:read
	- Install to Workspace
	- Bot User OAuth Token
		- Make a note of the token that begins with xoxb-.
- Event Subscriptions
	- Enable Events
		- On
	- Subscribe to bot events
	- Add Bot User Event
		- message.channels
	- Save Changes

2. on Slack App
	- invite bot
		- @(Name)
	- invite

note) Have them participate in all the channels where you want to collect incidents.

3. your os terminal
	- set environment
		- windows
			- set SLACK_APP_TOKEN=xapp-...
			- set SLACK_BOT_TOKEN=xoxb-...
		- linux
			- export SLACK_APP_TOKEN=xapp-...
			- export SLACK_BOT_TOKEN=xoxb-...
	- run this tool

# usecase

1. Decide with your team what reaction you want the responded to be. -> config [Label]
2. Get the ID of the slack channel where you want the string to be detected. -> config 
3. Get the ID of the slack channel you want to have incidents checked. (Separate channels, etc. If necessary.)
4. Use the test option to configure the string and time you want to detect. -> config [Rules]
5. 


# config file

Configs format is tab split values. The definition is ignored if you put sharp(#) at the beginning.

## auto read suppot

config file supported auto read. so, you rewrite config file, tool not necessaly rerun, tool just read this.

## [Rules]

Define rules for detecting messages.

```
[Rules]
.*Error.*	2[0-9]:.*:.*	CHANNEL1
```

1. include message (can use regex.)
2. date range (can use regex.)
3. channel label. If detect rule, post message to channel.

note) Date Format is "2006/01/02 15:04:05 Mon(-Sun)".
  If you detect message include "Fault" and every day at 10:00-12:00, rule is
  ```
  .*Fault.* .*/.*/.* 1[0-2]:.*:.* .*
  ```

note) not only single but can write plural rules by tsv.

## [Incidents]

This config for incidents managed channel.

```
[Incidents]
CHANNEL1	C021C9EC76C	20	[Rule1]	10
```

1. channel label
2. channnel id for Incidents
3. post message used this to header. (you use to analyze label.)
4. Number of message to go back

note) 4. is too big, check more slowly..<br>

note) not only single but can write plural rules by tsv.

## [Label]

Define which reactions are marked as supported.

note) [This page is a good reference for what marks can be used.](https://qiita.com/yamadashy/items/ae673f2bae8f1525b6af)

## example

```
[Rules]
.*Error.*	.*:.*:.*	CHANNEL1
[Incidents]
CHANNEL1	C021C9EC76C	[Rule1]	10
[Label]
white_check_mark
```

# options

```
  -auto
        [-auto=config auto read/write mode (true is enable)] (default true)
  -config string
        [-config=config file)] (default "IMS.ini")
  -debug
        [-debug=debug mode (true is enable)]
  -log
        [-log=logging mode (true is enable)]
  -loop int
        [-loop=incident check loop time. ] (default 30)
  -onlyincident
        [-onlyincident=incident check and exit mode.]
  -reverse
        [-reverse=check rule to reverse (true is enable)]
  -test string
        [-test=Test what happens when you set the message.]
  -verbose
        [-verbose=check output verbose (true is enable)] (default true)
```

## -auto

## -config

Specify the configuration file.

## -debug

Run in the mode that outputs various logs.

## -log

Specify the log file name.

## -loop

Interval between incidents checks (in seconds)

## -onlyincident

If this option is specified, the tool will exit after the incidents check.

note) This can be used if you want to move the incident manually.

## -reverse

all check rules to reverse

## -test

Give a message to check the rule and exit the tool.

```
>IMS -test="Error test"

[Test] Error test
this message include rule (1)!
```

note) The number in parentheses () indicates the number of rules that have been matched.

## -verbose

Displays not only unsupported messages, but also supported ones.

```
[message] norml message [date] 2021/05/02 21:27:28
[message] test message [date] 2021/05/02 17:54:05
```

to

```
NG [message] norml message [date] 2021/05/02 21:27:28
OK [message] error and reactioned [date] 2021/05/02 21:25:40 [user]  yasutakato
NG [message] test message [date] 2021/05/02 17:54:05
```

note) In this mode, the name of the person who responded will also be displayed.

# license

Apache-2.0 License<br>
BSD-2-Clause License<br>
BSD-3-Clause License
