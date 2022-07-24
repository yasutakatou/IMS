# IMS

**Incident management tool with slack**.

### v0.2

- add hotline mode
	-  Mention a specific message to make it an incident.
- add check target id
	- You can specify the ID to check for messages.

### v0.3

- It can now be defined by name or channel name instead of ID.
	- Pre) U024ZT3BHU5	After) adminuser

### v0.4

- add Reacji Channeler mode.
	- Added support for visualized messages using **Reacji Channeler**.

### 20220409

Add how to respond to **Private channel**.<br>
<br>
note) No change in code.<br>

### v0.5

- I can now report mixed A messages.

![image](https://user-images.githubusercontent.com/22161385/164874744-97e2027d-4554-4d2e-bc97-47eed8932f75.png)

### v0.6

- Added **reminder function** and ability to **delete remarks on reminder channels**.

![image](https://user-images.githubusercontent.com/22161385/164969202-53e05b2a-631e-4e4d-b673-8ac2365cd1e8.png)

### v0.7

- Fixed a bug in Reacji mode that prevented the system from responding to text-only messages.

### v0.8

- Allow defining user IDs to be forwarded by **Reacji**.

### v0.9

- Added the ability to **check the speaker** in the incident management channel.

### v0.91

- Alerts with referring URLs now forward content to the default channel
	- With Reaacji Channeler , it is difficult to understand the content of forwarded alerts, so we have made the forwarding at least include the content of those with links.

# Solution

As the center of communication at work has been replaced from e-mail to chat, you may have changed the alert notification destination of your monitoring tool to chat.
But hasn't it changed that people are reading the **all messages and making decisions and responses**?<br>
Let's change that!<br>
This tool enables easy incident management through chat, accelerating ChatOps!!<br>

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
			- v0.3 If you want to use **-idlookup** mode, you also need to define the following
				- channels:read
				- groups:read
				- im:read
				- mpim:read
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

## 20220409 

If you want to use **Private channnel**, add the following settings

- OAuth & Permissions
	- Scopes
	- Bot Token Scopes
		- **groups:history**
	- Install to Workspace
- Event Subscriptions
	- Subscribe to bot events
	- Add Bot User Event
		- **message.groups**
	- Save Changes

3. your OS terminal
	- set environment
		- windows
			- set SLACK_APP_TOKEN=xapp-...
			- set SLACK_BOT_TOKEN=xoxb-...
		- linux
			- export SLACK_APP_TOKEN=xapp-...
			- export SLACK_BOT_TOKEN=xoxb-...
	- run this tool

## v0.4) set up for Reacji Channeler

about **Reacji Channeler**

[Reacji Channeler](https://reacji-channeler.builtbyslack.com/)<br>
[Slack 用リアク字チャンネラー](https://slack.com/intl/ja-jp/help/articles/360000482666-Slack-%E7%94%A8%E3%83%AA%E3%82%A2%E3%82%AF%E5%AD%97%E3%83%81%E3%83%A3%E3%83%B3%E3%83%8D%E3%83%A9%E3%83%BC)

Define a forwarding reaction for the channel that **collects the incidents**.

![image](https://user-images.githubusercontent.com/22161385/142754689-e65a1b1e-5c2b-4505-8d48-1f933be7ddb9.png)

If the rule is met, it will be automatically **marked and forwarded**.

![3](https://user-images.githubusercontent.com/22161385/142754729-4aa99751-d11b-4d94-b0f8-248c9aa3033a.png)<br>
![4](https://user-images.githubusercontent.com/22161385/142754749-20ea0ba9-daa4-4fa2-b3bd-1c2f4841f75d.png)

**Reports with links** will go up on the channel for reporting.

![image](https://user-images.githubusercontent.com/22161385/142754760-225ed01a-00b1-48e0-b65d-3d38ff98b2d2.png)

The mode of **-reverse** is also supported

![image](https://user-images.githubusercontent.com/22161385/142754781-b8eb6394-46be-4ca8-ac3f-cd1837d60e73.png)<br>
![image](https://user-images.githubusercontent.com/22161385/142754787-068bc919-d9a0-41e8-9bcb-0f4ca0723563.png)<br>
![image](https://user-images.githubusercontent.com/22161385/142754790-49a72621-0f04-4a1d-a1be-7f2ba55ac850.png)<br>

note) If you are in mode A, **you will not be able to add tags** to your report.

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

note) **not only single define but can write plural rules**.<br>

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
note) **not only single define but can write plural rules**.<br>
note) v0.3: You can also specify a channel name instead of an ID.<br>

### Special Definition

In the case of **-reverse mode**, it defines the **default incident registration destination** when all the rules are not match.

```
DEFAULT	C025FKF3QJV	[Alert]
```

1. **"DEFAULT"** is static define.
2. channnel id for Incident manage.
3. message is use this header.

note) v0.3: You can also specify a channel name instead of an ID.<br>

## [Label]

Define which reactions are **marked as resolved**.<br>
<br>
note) [This page is a good reference for what marks can be used.](https://qiita.com/yamadashy/items/ae673f2bae8f1525b6af)

## [Report]

Define the channel for **report output**.<br>
The default **cycle is once a day**, but you can change it with option -loop.<br>

note) v0.3: You can also specify a channel name instead of an ID.<br>

## [PostID]

Messages from the ID defined here will be **checked**.

note) **not only single define but can write plural IDs**.<br>
note) You can also specify the ID of the **bot**.<br>
note) v0.3: You can also specify a user name instead of an ID.<br>

![image](https://user-images.githubusercontent.com/22161385/125152354-fe1b4780-e186-11eb-9cd0-8f33940ce16a.png)

## [Hotline]

Defines the destination for mailed incidents.<br>

```
Hot1	U024ZT3BHU5	here
```

1. label for define.
2-. mention ids

note) **not only single define but can write plural rules**.<br>
note) Slack user ID or **here, channnel, everyone** can be defined.<br>
note) v0.3: You can also specify a user name instead of an ID.<br>

## [Reacji]

Define which reactions for **Reacji Channeler**.<br>
Forward the incident to the channel that collects it with this definition.<br>
<br>
note) [This page is a good reference for what marks can be used.](https://qiita.com/yamadashy/items/ae673f2bae8f1525b6af)

```
warning
```

## [Reminder]

This function periodically **picks up unaddressed incidents** and notifies you of them.<br>
Set the channel and time to be notified.<br>

note) The first part is the channel name.<br>
note) Specify the time you want to be notified in **tab-delimited** format using a **regular expression**.<br>

```
alert	.*1.*	.*2.*
````

note) In the above example, We'll keep you posted on Channel A from 10-24.<br>
note) **not only single define but can write plural rules**.<br>

## [ReacjiID]

The user ID defined here will be **transferred by Reacji**. Specifies primarily **webhook bots**.<br>

```
datadog
```


note) **not only single define but can write plural rules**.<br>

## [MgmtReport]

Added the ability to **check the speaker** in the incident management channel.<br>
Messages from the ID defined here will be **checked** and output to channel for **report**.<br>
**If empty, all submissions are forwarded to the reporting channel.**

```
[MgmtReport]
U024ZT3BHU5
```

Similar to [PostID], but this function was implemented because there was a request that it would be better to handle **manager's instructions as Issues** when used by a team.<br>
<br>
note) **not only single define but can write plural IDs**.<br>
note) You can also specify the ID of the **bot**.<br>
note) You can also specify a user name instead of an ID.<br>

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

### v0.3

```
[Rules]
.*Error.*	.*:.*:.*	[RuleX]	CHANNEL1	Hot1
.*Warn.*	.*:.*:.*	[RuleX]	CHANNEL1	No
.*Info.*	.*:.*:.*	[RuleX]	CHANNEL1	No
.*Debug.*	.*:.*:.*	[RuleX]	CHANNEL1	No
[Incidents]
CHANNEL1	incidents	20
DEFAULT	incidents	[Alert]
[Label]
white_check_mark
[Report]
report
[PostID]
ims
adminuser
[Hotline]
Hot1	adminuser	here
```

### v0.4

```
[Rules]
.*Error.*	.*:.*:.*	[Error]	CHANNEL1	Hot1
.*Warn.*	.*:.*:.*	[Warn]	CHANNEL1	No
.*Info.*	.*:.*:.*	[Info]	CHANNEL1	No
.*Debug.*	.*:.*:.*	[Debug]	CHANNEL1	No
[Incidents]
CHANNEL1	incidents	20
DEFAULT	incidents	[Alert]	
[Label]
white_check_mark
[Report]
rep
[PostID]
user
adminuser
[Hotline]
Hot1	adminuser	here
[Reacji]
warning
```

### v0.6

```
[Rules]
.*Error.*	.*:.*:.*	[Error]	CHANNEL1	Hot1
.*Warn.*	.*:.*:.*	[Warn]	CHANNEL1	No
.*Info.*	.*:.*:.*	[Info]	CHANNEL1	No
.*Debug.*	.*:.*:.*	[Debug]	CHANNEL1	No
[Incidents]
CHANNEL1	incidents	20
DEFAULT	incidents	[Alert]	
[Label]
white_check_mark
[Report]
rep
[PostID]
user
adminuser
[Hotline]
Hot1	adminuser	here
[Reacji]
warning
[Reminder]
alert	.*1.*	.*2.*
```

### v0.8

```
[Rules]
.*Error.*	.*:.*:.*	[Error]	CHANNEL1	Hot1
.*Warn.*	.*:.*:.*	[Warn]	CHANNEL1	No
.*Info.*	.*:.*:.*	[Info]	CHANNEL1	No
.*Debug.*	.*:.*:.*	[Debug]	CHANNEL1	No
[Incidents]
CHANNEL1	incidents	20
DEFAULT	incidents	[Alert]	
[Label]
white_check_mark
[Report]
rep
[PostID]
user
adminuser
[Hotline]
Hot1	adminuser	here
[Reacji]
warning
[Reminder]
alert	.*1.*	.*2.*
[ReacjiID]
datadog
```

### v0.9

```
[Rules]
.*Error.*	.*:.*:.*	[Error]	CHANNEL1	Hot1
.*Warn.*	.*:.*:.*	[Warn]	CHANNEL1	No
.*Info.*	.*:.*:.*	[Info]	CHANNEL1	No
.*Debug.*	.*:.*:.*	[Debug]	CHANNEL1	No
[Incidents]
CHANNEL1	incidents	20
DEFAULT	incidents	[Alert]	
[Label]
white_check_mark
[Report]
rep
[PostID]
user
adminuser
[Hotline]
Hot1	adminuser	here
[Reacji]
warning
[Reminder]
alert	.*1.*	.*2.*
[ReacjiID]
datadog
[MgmtReport]
```

# options

```
  -auto
        [-auto=config auto read/write mode (true is enable)] (default true)
  -clearReminder
        [-clearReminder=clear reminder channel and exit mode.]
  -config string
        [-config=config file)] (default "IMS.ini")
  -debug
        [-debug=debug mode (true is enable)]
  -idlookup
        [-idlookup=resolve to ID definition (true is enable)] (default true)
  -log
        [-log=logging mode (true is enable)]
  -loop int
        [-loop=incident check loop time (Hour). ] (default 24)
  -onlyReport
        [-onlyReport=incident check and exit mode.]
  -reacji
        [-reacji=Slack: reacji channeler mode (true is enable)]
  -reminder int
        [-reminder=Interval for posting reminders (Seconds). ] (default 30)
  -reverse
        [-reverse=check rule to reverse (true is enable)]
  -test string
        [-test=Test what happens when you set the message.]
  -verbose
        [-verbose=incident output verbose (true is enable)]
```

## -auto

config auto read/write mode.

## -clearReminder

Turn off messages in the Reminders channel.

## -config

Specify the configuration file name.

## -debug

Run in the mode that outputs various logs.

## -idlookup

When enabled, it converts channel names and user names into IDs.

## -log

Specify the log file name.

## -loop

Interval between incidents checks (in Hours). **Default is 24 Hour**.

## -onlyReport

If this option is specified, the tool will exit after the incidents report.<br>

note) This can be used if you want to check the incident report manually.<br>

## -reacji

Activate the mode that uses **Reacji Channeler**.

## -reminder int

The interval at which to check for reminders.

note) Units are in **seconds**.<br>

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
