# IMS

**Incident management tool with slack**.

### v0.3

- add hotline mode
	-  Mention a specific message to make it an incident.
- add check target id
	- You can specify the ID to check for messages.

# Solution

As the center of communication at work has been replaced from e-mail to chat, you may have changed the alert notification destination of your monitoring tool to chat.
But hasn't it changed that people are reading the **all messages and making decisions and responses**?
Let's change that!
This tool enables easy incident management through chat, accelerating ChatOps!!

# Feature

This tool has three major functions.<br>

1. check the posted messages according to the rules, label them if they fit the rules, and repost them to the incidents management channel.

![2](https://user-images.githubusercontent.com/22161385/117542584-2f2aaf00-b054-11eb-8405-570a21570f9c.png)

It is also possible to post actions that do not fit the reverse rule

2. Check if there are any reactions in the incidents management channel, and output the ones that are not there.

![1](https://user-images.githubusercontent.com/22161385/117542581-2d60eb80-b054-11eb-893e-d61355040935.png)

3. Periodically post a list of unsupported alerts to report channel.

![image](https://user-images.githubusercontent.com/22161385/122943521-10893900-d3b2-11eb-8d90-db3896fa2d6b.png)

In other words, you can use the following

1. Invite this tool to the channel you are throwing the monitoring message into. The tool checks in the all messages
2. Periodically, the report will run and display a list of messages that have not been reacted to, so check the unacted ones and leave a history of your teams action in the thread.

This makes it possible to

1. **Identify unanswered alerts**
2. **Filter for known messages**
3. **Your team can keep a history of responses to alerts.**

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
	- Generate Token and Scopes
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
		- files:write
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

note) Bot have them participate in all the channels where you want to collect incidents.

3. your OS terminal
	- set environment
		- windows
			- set SLACK_APP_TOKEN=xapp-...
			- set SLACK_BOT_TOKEN=xoxb-...
		- linux
			- export SLACK_APP_TOKEN=xapp-...
			- export SLACK_BOT_TOKEN=xoxb-...
	- run this tool

# usecase

1. What alert messages will you respond to? Decide with your team what alert messages you will respond to, or ignore. -> config [Label]
2. Decide on a channel for message retrieval, a channel for incident management, and a channel for reporting. -> config [Incidents]
3. ecide which reaction mark will be used to mark the item as handled. -> config [Label]
4. Define the channel for report output.-> config [Report]

# config file

Configs format is **tab split values**. The definition is ignore if you put sharp(#) at the beginning.

## auto read suppot

config file supported **auto read. so**, you rewrite config file, **tool not necessaly rerun, tool just this**.

## [Rules]

Define rules for detecting messages.

```
[Rules]
.*Error.*	.*:.*:.*	[RuleX]	CHANNEL1	Hot1
```

1. strings define (can use regex.) note: The meaning of the string to be included.
2. Date and time range (can use regex.)
3. Give this label to messages that match the rule. (you use to analyze messages.)
4. channel label. If detect rule, post message to channel defined.
5. Mention it and make it an incident. Define the name of the **[Hotline] label**.

note) Date Format is "2006/01/02 15:04:05 Mon(-Sun)".<br>
  If you want to detect message include "Fault" and every day at 10:00-12:00, rule is<br>
  ```
  .*Fault.* .*/.*/.* 1[0-2]:.*:.* .*
  ```

note) **not only single define but can write plural rules**.

## [Incidents]

This config for incidents managed channel.

```
[Incidents]
CHANNEL1	C025FKF3QJV	20
```

1. label for channel.
2. channnel id for Incident manage.
3. Number of message to go back reference.

note) 3. is too big, check more slowly..<br>
note) **not only single define but can write plural rules**.

### Special Definition

In the case of **-reverse mode**, it defines the **default incident registration destination** when all the rules are not match.

```
DEFAULT	C025FKF3QJV	[Alert]
```

1. **"DEFAULT"** is static define.
2. channnel id for Incident manage.
3. message is use this header.

## [Label]

Define which reactions are **marked as resolved**.

note) [This page is a good reference for what marks can be used.](https://qiita.com/yamadashy/items/ae673f2bae8f1525b6af)

## [Report]

Define the channel for **report output**.<br>
The default **cycle is once a day**, but you can change it with option -loop.

## [PostID]

Messages from the ID defined here will be **checked**.

note) **not only single define but can write plural IDs**.
note) You can also specify the ID of the **bot**.

![image](https://user-images.githubusercontent.com/22161385/125152354-fe1b4780-e186-11eb-9cd0-8f33940ce16a.png)

## [Hotline]

Defines the destination for mailed incidents.

```
Hot1	U024ZT3BHU5	here
```

1. label for define.
2-. mention ids

note) **not only single define but can write plural rules**.
note) Slack user ID or **here, channnel, everyone** can be defined.

## example

```
[Rules]
.*Error.*	.*:.*:.*	[RuleX]	CHANNEL1	Hot1
.*Warn.*	.*:.*:.*	[RuleX]	CHANNEL1	No
[Incidents]
CHANNEL1	C025FKF3QJV	20
DEFAULT	C025FKF3QJV	[Alert]
[Label]
white_check_mark
[Report]
C0256BTKP54
[PostID]
U024ZT3BHU5
[Hotline]
Hot1	U024ZT3BHU5	here
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
        [-loop=incident check loop time. ] (default 24)
  -onlyReport
        [-onlyReport=incident check and exit mode.]
  -reverse
        [-reverse=check rule to reverse (true is enable)]
  -test string
        [-test=Test what happens when you set the message.]
  -verbose
        [-verbose=check output verbose (true is enable)] (default true)
```

## -auto

config auto read/write mode.

## -config

Specify the configuration file name.

## -debug

Run in the mode that outputs various logs.

## -log

Specify the log file name.

## -loop

Interval between incidents checks (in Hours). **Default is 24 Hour**.

## -onlyReport

If this option is specified, the tool will exit after the incidents report.<br>

note) This can be used if you want to check the incident report manually.

## -reverse

**all check rules to reverse**

![1](https://user-images.githubusercontent.com/22161385/122678154-75148e80-d220-11eb-9c71-873538f8b5ed.png)

It will be reversed as follows

![2](https://user-images.githubusercontent.com/22161385/122678156-7776e880-d220-11eb-8fbe-d7e86222baca.png)

note) Rules in hotline mode will be made **incidental even if they are reversed**.

## -test

If this option is specified, the tool will exit after the message check.<br>
This can be used if you want to check the message check manually.
 
```
>IMS -test="Error test"

[Test] Error test
this message include rule (1)!
```

note) The number in parentheses () indicates **the order of rules that have been matched**.

## -verbose

Displays not only unsolved messages, but also solved ones.

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

note) In this mode, **the name of the person who resolve will also be displayed**.

# license

Apache-2.0 License<br>
BSD-2-Clause License<br>
BSD-3-Clause License
