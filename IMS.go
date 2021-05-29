/*
 * Incident management tool with slack.
 *
 * @author    yasutakatou
 * @copyright 2021 yasutakatou
 * @license   Apache-2.0 License, BSD-2 Clause License, BSD-3 Clause License
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/slack-go/slack/socketmode"
	"gopkg.in/ini.v1"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type ruleData struct {
	TARGET  string
	EXCLUDE string
	LABEL   string
}

type incidentData struct {
	LABEL    string
	CHANNNEL string
	HEAD     string
	LIMIT    int
}

var (
	debug, logging bool
	label          string
	incidents      []incidentData
	rules          []ruleData
)

func main() {
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Logging := flag.Bool("log", false, "[-log=logging mode (true is enable)]")
	_Config := flag.String("config", "IMS.ini", "[-config=config file)]")
	_loop := flag.Int("loop", 30, "[-loop=incident check loop time. ]")
	_onlyincident := flag.Bool("onlyincident", false, "[-onlyincident=incident check and exit mode.]")
	_verbose := flag.Bool("verbose", false, "[-verbose=incident output verbose (true is enable)]")
	_test := flag.String("test", "", "[-test=Test what happens when you set the message.]")
	_autoRW := flag.Bool("auto", true, "[-auto=config auto read/write mode (true is enable)]")
	_reverse := flag.Bool("reverse", false, "[-reverse=check rule to reverse (true is enable)]")

	flag.Parse()

	debug = bool(*_Debug)
	logging = bool(*_Logging)

	if Exists(*_Config) == true {
		loadConfig(*_Config)
	} else {
		fmt.Printf("Fail to read config file: %v\n", *_Config)
		os.Exit(1)
	}

	if *_test != "" {
		testRule(*_test, *_reverse)
		os.Exit(0)
	}

	if *_onlyincident == true {
		incident(*_verbose)
		os.Exit(0)
	}

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	if *_autoRW == true {
		go func() {
			for {
				select {
				case <-watcher.Events:
					loadConfig(*_Config)
				case <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()
	}

	if err := watcher.Add(*_Config); err != nil {
		fmt.Println("ERROR", err)
	}

	ruleChecker(*_reverse)

	for {
		incident(*_verbose)
		time.Sleep(time.Second * time.Duration(*_loop))
	}
	os.Exit(0)
}

func testRule(message string, reverse bool) {
	fmt.Println("[Test] " + message)

	result := checkMessage(message, reverse)
	if result != 0 {
		fmt.Printf("this message include rule (%d)!\n", result)
	} else {
		fmt.Println("this message exclude rules..")
	}
}

func incident(verbose bool) {
	const layout = "2006/01/02 15:04:05"
	t := time.Now()
	fmt.Println(" - - " + t.Format(layout) + " - -")

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	api := slack.New(botToken)

	for i := 0; i < len(incidents); i++ {
		debugLog("incident: " + incidents[i].CHANNNEL)
		params := slack.GetConversationHistoryParameters{ChannelID: incidents[i].CHANNNEL, Limit: incidents[i].LIMIT}
		messages, err := api.GetConversationHistory(&params)
		if err != nil {
			fmt.Printf("incident not get: %s\n", err)
			return
		}
		for _, message := range messages.Messages {
			name := checkReaction(message.Reactions)
			if verbose == true {
				if name == "" {
					fmt.Println("NG [message] " + message.Text + " [date] " + convertTime(message.Timestamp))
				} else {
					fmt.Println("OK [message] " + message.Text + " [date] " + convertTime(message.Timestamp) + " [user] " + name)
				}
			} else {
				if name == "" {
					fmt.Println("[message] " + message.Text + " [date] " + convertTime(message.Timestamp))
				}
			}
		}
	}
}

func convertTime(unixTime string) string {
	var tsStr string
	if strings.Index(unixTime, ".") != -1 {
		tss := strings.Split(unixTime, ".")
		tsStr = tss[0]
	} else {
		tsStr = unixTime
	}
	ts, _ := strconv.ParseInt(tsStr, 10, 64)
	t := time.Unix(ts, 0)
	const layout = "2006/01/02 15:04:05"
	return t.Format(layout)
}

func checkReaction(reactions []slack.ItemReaction) string {
	for _, reaction := range reactions {
		if reaction.Name == label {
			users := ""
			for _, user := range reaction.Users {
				users = users + " " + getUsername(user)
			}
			return users
		}
	}
	return ""
}

func getUsername(userID string) string {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	user, err := api.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}
	return user.Profile.RealName
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func loadConfig(configFile string) {
	loadOptions := ini.LoadOptions{}
	loadOptions.UnparseableSections = []string{"Rules", "Incidents", "Label"}

	rules = nil
	incidents = nil
	label = ""

	cfg, err := ini.LoadSources(loadOptions, configFile)
	if err != nil {
		fmt.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}

	setStructs("Rules", cfg.Section("Rules").Body(), 0)
	setStructs("Incidents", cfg.Section("Incidents").Body(), 1)
	setStructs("Label", cfg.Section("Label").Body(), 2)
}

func setStructs(configType, datas string, flag int) {
	debugLog(" -- " + configType + " --")

	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 && flag != 2 {
			if strings.Index(v, "\t") != -1 {
				strs := strings.Split(v, "\t")

				switch flag {
				case 0:
					if len(strs) == 3 {
						rules = append(rules, ruleData{TARGET: strs[0], EXCLUDE: strs[1], LABEL: strs[2]})
						debugLog(v)
					}
				case 1:
					if len(strs) == 4 {
						convInt, err := strconv.Atoi(strs[3])
						if err == nil {
							incidents = append(incidents, incidentData{LABEL: strs[0], CHANNNEL: strs[1], HEAD: strs[2], LIMIT: convInt})
							debugLog(v)
						}
					}
				}
			}
		} else if flag == 2 {
			label = v
			debugLog(v)
		}
	}
}

func debugLog(message string) {
	var file *os.File
	var err error

	if debug == true {
		fmt.Println(message)
	}

	if logging == false {
		return
	}

	const layout = "2006-01-02_15"
	t := time.Now()
	filename := "inco_" + t.Format(layout) + ".log"

	if Exists(filename) == true {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	fmt.Fprintln(file, message)
}

func postMessage(channelInt int, message string) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	_, _, err := api.PostMessage(incidents[channelInt].CHANNNEL, slack.MsgOptionText(incidents[channelInt].HEAD+" "+message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func ruleChecker(reverse bool) {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(debug),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}
				if eventsAPIEvent.Type == slackevents.CallbackEvent {
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						debugLog("receive message: " + ev.Text)
						result := checkMessage(ev.Text, reverse)
						if result != 0 && channelMatch(ev.Channel) == false {
							postMessage(result-1, ev.Text)
						} else if channelMatch(ev.Channel) == false {
							markReaction(ev.Channel, ev.TimeStamp)
						}
					}
				}
				client.Ack(*evt.Request)
			}
		}
	}()

	go client.Run()
}

func channelMatch(channel string) bool {
	for i := 0; i < len(incidents); i++ {
		if incidents[i].CHANNNEL == channel {
			return true
		}
	}
	return false
}

func markReaction(channnel, ts string) {
	msgRef := slack.NewRefToMessage(channnel, ts)

	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	if err := api.AddReaction(label, msgRef); err != nil {
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	}
}

func checkMessage(message string, reverse bool) int {
	wdays := [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	const layout = "2006/01/02 15:04:05"
	t := time.Now()
	nowDate := t.Format(layout) + " " + wdays[t.Weekday()]

	for i := 0; i < len(rules); i++ {
		debugLog("messageRegex: " + rules[i].TARGET)
		messageRegex := regexp.MustCompile(rules[i].TARGET)

		if reverse == true {
			if messageRegex.MatchString(message) == false {
				debugLog("messageRegex: ok")
				debugLog("nowDate: " + nowDate)

				debugLog("dateRegex: " + rules[i].EXCLUDE)
				dateRegex := regexp.MustCompile(rules[i].EXCLUDE)
				if dateRegex.MatchString(nowDate) == false {
					debugLog("dateRegex: ok")
					if act := incidentCheck(rules[i].LABEL); act != 0 {
						return act
					}
				}
			}
		} else {
			if messageRegex.MatchString(message) == true {
				debugLog("messageRegex: ok")
				debugLog("nowDate: " + nowDate)

				debugLog("dateRegex: " + rules[i].EXCLUDE)
				dateRegex := regexp.MustCompile(rules[i].EXCLUDE)
				if dateRegex.MatchString(nowDate) == true {
					debugLog("dateRegex: ok")
					if act := incidentCheck(rules[i].LABEL); act != 0 {
						return act
					}
				}
			}
		}
	}
	return 0
}

func incidentCheck(incidentName string) int {
	for i := 0; i < len(incidents); i++ {
		if incidents[i].LABEL == incidentName {
			return i + 1
		}
	}
	return 0
}
