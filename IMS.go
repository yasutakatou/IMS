/*
 * Incident management tool with slack.
 *
 * @author    yasutakatou
 * @copyright 2021 yasutakatou
 * @license   Apache-2.0 License, BSD-2 Clause License, BSD-3 Clause License
 */
package main

import (
	"encoding/json"
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
	HEAD    string
	LABEL   string
	HOTLINE string
}

type incidentData struct {
	LABEL    string
	CHANNNEL string
	LIMIT    int
}

type alertData struct {
	LABEL string
	USERS []string
}

var (
	debug, logging, reacji bool
	label, reacjiStr       string
	defaultChannel         []string
	report                 string
	incidents              []incidentData
	rules                  []ruleData
	postids                []string
	alerts                 []alertData
)

func main() {
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Logging := flag.Bool("log", false, "[-log=logging mode (true is enable)]")
	_Config := flag.String("config", "IMS.ini", "[-config=config file)]")
	_loop := flag.Int("loop", 24, "[-loop=incident check loop time (Hour). ]")
	_onlyReport := flag.Bool("onlyReport", false, "[-onlyReport=incident check and exit mode.]")
	_verbose := flag.Bool("verbose", false, "[-verbose=incident output verbose (true is enable)]")
	_test := flag.String("test", "", "[-test=Test what happens when you set the message.]")
	_autoRW := flag.Bool("auto", true, "[-auto=config auto read/write mode (true is enable)]")
	_reverse := flag.Bool("reverse", false, "[-reverse=check rule to reverse (true is enable)]")
	_IDLookup := flag.Bool("idlookup", true, "[-idlookup=resolve to ID definition (true is enable)]")
	_reacji := flag.Bool("reacji", false, "[-reacji=Slack: reacji channeler mode (true is enable)]")

	flag.Parse()

	debug = bool(*_Debug)
	logging = bool(*_Logging)
	reacji = bool(*_reacji)

	if *_test != "" {
		testRule(*_test, *_reverse)
		os.Exit(0)
	}

	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must be set.\n")
		os.Exit(1)
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

	if Exists(*_Config) == true {
		loadConfig(api, *_Config, *_IDLookup)
	} else {
		fmt.Printf("Fail to read config file: %v\n", *_Config)
		os.Exit(1)
	}

	if *_onlyReport == true {
		incident(api, *_verbose, *_reverse)
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
					loadConfig(api, *_Config, *_IDLookup)
				case <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()
	}

	if err := watcher.Add(*_Config); err != nil {
		fmt.Println("ERROR", err)
	}

	ruleChecker(api, *_reverse)

	for {
		incident(api, *_verbose, *_reverse)
		time.Sleep(time.Hour * time.Duration(*_loop))
	}
	os.Exit(0)
}

func testRule(message string, reverse bool) {
	debugLog("[Test] " + message)

	result, ruleInt := checkMessage(message)
	if result != 0 {
		fmt.Printf("this message include rule (%d)!\n", ruleInt)
	} else {
		fmt.Println("this message exclude rules..")
	}
}

func incident(api *slack.Client, verbose, reverse bool) {
	const layout = "2006/01/02 15:04:05"
	t := time.Now()

	ret := ""
	dates := " - - " + t.Format(layout) + " - -"
	debugLog(ret)

	for i := 0; i < len(incidents); i++ {
		debugLog("incident: " + incidents[i].CHANNNEL)
		params := slack.GetConversationHistoryParameters{ChannelID: incidents[i].CHANNNEL, Limit: incidents[i].LIMIT}
		messages, err := api.GetConversationHistory(&params)
		if err != nil {
			fmt.Printf("incident not get: %s\n", err)
			return
		}
		for x, message := range messages.Messages {
			ret = ""
			if reacji == true {
				if x == 0 {
					postMessageStr(api, report, "", dates)
				}

				if reverse == true {
					name := checkReaction(api, message.Reactions)

					if strings.Index(message.Text, "Hotline Alert") == -1 {
						if verbose == true {
							if name == "" {
								stra := "NG [message] " + message.Text + " [date] " + convertTime(message.Timestamp)
								debugLog(stra)
								ret = ret + stra + "\n\n"
							} else {
								stra := "OK [message] " + message.Text + " [date] " + convertTime(message.Timestamp) + " [user] " + name
								debugLog(stra)
								ret = ret + stra + "\n\n"
							}
						} else {
							if name == "" {
								stra := "[message] " + message.Text + " [date] " + convertTime(message.Timestamp)
								debugLog(stra)
								ret = ret + stra + "\n\n"
							}
						}
						postMessageStr(api, report, "", ret)
					}
				} else {
					actualAttachmentJson, err := json.Marshal(message.Attachments)
					if err != nil {
						fmt.Println("expected no error unmarshaling attachment with blocks, got: %v", err)
					}
					mess := string(actualAttachmentJson)
					result, ruleInt := checkMessage(mess)
					name := checkReaction(api, message.Reactions)

					if result > 0 && strings.Index(message.Text, "Hotline Alert") == -1 {
						if verbose == true {
							if name == "" {
								stra := "NG [message] " + message.Text + " [date] " + convertTime(message.Timestamp)
								debugLog(stra)
								ret = ret + stra + "\n\n"
							} else {
								stra := "OK [message] " + message.Text + " [date] " + convertTime(message.Timestamp) + " [user] " + name
								debugLog(stra)
								ret = ret + stra + "\n\n"
							}
						} else {
							if name == "" {
								stra := "[message] " + message.Text + " [date] " + convertTime(message.Timestamp)
								debugLog(stra)
								ret = ret + stra + "\n\n"
							}
						}
						postMessageStr(api, report, rules[ruleInt].HEAD, ret)
					}
				}
			} else {
				mess := message.Text

				if len(mess) == 0 {
					actualAttachmentJson, err := json.Marshal(message.Attachments)
					if err != nil {
						fmt.Println("expected no error unmarshaling attachment with blocks, got: %v", err)
					}
					mess = string(actualAttachmentJson)
				}

				if len(mess) > 0 && mess != "null" {
					hlen := strings.Index(mess, "[Hotline Alert!]")
					if hlen != -1 {
						mess = mess[:hlen]
					}
					name := checkReaction(api, message.Reactions)

					if verbose == true {
						if name == "" {
							stra := "NG [message] " + mess + " [date] " + convertTime(message.Timestamp)
							debugLog(stra)
							ret = ret + stra + "\n\n"
						} else {
							stra := "OK [message] " + mess + " [date] " + convertTime(message.Timestamp) + " [user] " + name
							debugLog(stra)
							ret = ret + stra + "\n\n"
						}
					} else {
						if name == "" {
							stra := "[message] " + mess + " [date] " + convertTime(message.Timestamp)
							debugLog(stra)
							ret = ret + stra + "\n\n"
						}
					}
				}
				postTextFile(api, ret, report, dates)
			}
		}
	}

}

func postTextFile(api *slack.Client, strs, repChan, dates string) {
	params := slack.FileUploadParameters{
		Title:    dates,
		Filetype: "txt",
		Content:  strs,
		Channels: []string{repChan},
	}
	_, err := api.UploadFile(params)
	if err != nil {
		debugLog(fmt.Sprintf("%s\n", err))
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

func checkReaction(api *slack.Client, reactions []slack.ItemReaction) string {
	for _, reaction := range reactions {
		if reaction.Name == label {
			users := ""
			for _, user := range reaction.Users {
				users = users + " " + getUsername(api, user)
			}
			return users
		}
	}
	return ""
}

func getUsername(api *slack.Client, userID string) string {
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

func loadConfig(api *slack.Client, configFile string, IDLookup bool) {
	loadOptions := ini.LoadOptions{}
	loadOptions.UnparseableSections = []string{"Rules", "Incidents", "Label", "Report", "PostID", "Hotline", "Reacji"}

	rules = nil
	incidents = nil
	label = ""
	report = ""
	postids = nil
	alerts = nil

	cfg, err := ini.LoadSources(loadOptions, configFile)
	if err != nil {
		fmt.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}

	usersMap := map[string]string{}
	channelsMap := map[string]string{}

	if IDLookup == true {
		users, err := api.GetUsers()
		if err == nil {
			for _, user := range users {
				debugLog("UserIDs: " + user.ID + " " + user.Name)
				usersMap[user.Name] = user.ID
			}
		}
	}

	var cursor string
	for {
		requestParam := &slack.GetConversationsParameters{
			Types:           []string{"public_channel"},
			Limit:           1000,
			ExcludeArchived: false,
		}
		if cursor != "" {
			requestParam.Cursor = cursor
		}
		var channels []slack.Channel
		channels, cursor, err := api.GetConversations(requestParam)
		if err == nil {
			for _, channel := range channels {
				debugLog("ChannelIDs: " + channel.ID + " " + channel.Name)
				channelsMap[channel.Name] = channel.ID
			}
		}
		if cursor == "" {
			break
		}
	}

	setStructs(IDLookup, usersMap, channelsMap, "Rules", cfg.Section("Rules").Body(), 0)
	setStructs(IDLookup, usersMap, channelsMap, "Incidents", cfg.Section("Incidents").Body(), 1)
	setStructs(IDLookup, usersMap, channelsMap, "Label", cfg.Section("Label").Body(), 2)
	setStructs(IDLookup, usersMap, channelsMap, "Report", cfg.Section("Report").Body(), 3)
	setStructs(IDLookup, usersMap, channelsMap, "PostID", cfg.Section("PostID").Body(), 4)
	setStructs(IDLookup, usersMap, channelsMap, "Hotline", cfg.Section("Hotline").Body(), 5)
	setStructs(IDLookup, usersMap, channelsMap, "Reacji", cfg.Section("Reacji").Body(), 6)
}

func setStructs(IDLookup bool, users, channels map[string]string, configType, datas string, flag int) {
	debugLog(" -- " + configType + " --")

	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 && flag != 2 && flag != 3 && flag != 4 && flag != 6 {
			if strings.Index(v, "\t") != -1 {
				strs := strings.Split(v, "\t")

				switch flag {
				case 5:
					if len(strs) > 1 {
						var strr []string

						for i := 1; i < len(strs); i++ {
							strr = append(strr, setUserStr(IDLookup, users, strs[i]))
						}
						alerts = append(alerts, alertData{LABEL: strs[0], USERS: strr})
					}
				case 1:
					if strs[0] == "DEFAULT" {
						defaultChannel = append(defaultChannel, setChannelStr(IDLookup, channels, strs[1]))
						defaultChannel = append(defaultChannel, strs[2])
						debugLog("default channel: " + strs[0] + " " + setChannelStr(IDLookup, channels, strs[1]) + " " + strs[2])
					} else if len(strs) == 3 {
						convInt, err := strconv.Atoi(strs[2])
						if err == nil {
							incidents = append(incidents, incidentData{LABEL: strs[0], CHANNNEL: setChannelStr(IDLookup, channels, strs[1]), LIMIT: convInt})
							debugLog("add channel: " + strs[0] + " " + setChannelStr(IDLookup, channels, strs[1]) + " " + strs[2])
						}
					}
				case 0:
					if len(strs) == 5 {
						rules = append(rules, ruleData{TARGET: strs[0], EXCLUDE: strs[1], HEAD: strs[2], LABEL: strs[3], HOTLINE: strs[4]})
						debugLog(v)
					}
				}
			}
		} else if flag == 2 {
			label = v
			debugLog(v)
		} else if flag == 3 {
			report = setChannelStr(IDLookup, channels, v)
			debugLog(v)
		} else if flag == 4 {
			postids = append(postids, setUserStr(IDLookup, users, v))
			debugLog(v)
		} else if flag == 6 {
			reacjiStr = v
			debugLog(v)
		}
	}
}

func setChannelStr(IDLookup bool, channels map[string]string, key string) string {
	if IDLookup == true {
		us, ok := channels[key]
		if ok == true {
			debugLog("Resove Channels: " + key + " -> " + us)
			return us
		}
	}
	debugLog("No Resolv:" + key)
	return key
}

func setUserStr(IDLookup bool, users map[string]string, key string) string {
	if IDLookup == true {
		us, ok := users[key]
		if ok == true {
			debugLog("Resove User: " + key + " -> " + us)
			return us
		}
	}
	debugLog("No Resolv:" + key)
	return key
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

func postMessage(api *slack.Client, channelInt, ruleInt int, message string) {
	if len(message) < 5 {
		return
	}
	debugLog("POST channnel: " + incidents[channelInt].CHANNNEL + " label: " + rules[ruleInt].HEAD + " mess: " + message)
	_, _, err := api.PostMessage(incidents[channelInt].CHANNNEL, slack.MsgOptionText(rules[ruleInt].HEAD+" "+message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func postMessageStr(api *slack.Client, channelStr, channelLabel string, message string) {
	if len(message) < 5 {
		return
	}
	debugLog("POST channnel: " + channelStr + " label: " + channelLabel + " mess: " + message)
	_, _, err := api.PostMessage(channelStr, slack.MsgOptionText(channelLabel+" "+message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func ruleChecker(api *slack.Client, reverse bool) {
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
					fmt.Printf("evt!")
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						postId := ""
						if len(ev.BotID) > 0 {
							postId = ev.BotID
						} else {
							postId = ev.User
						}

						mess := ev.Text

						if len(mess) == 0 {
							actualAttachmentJson, err := json.Marshal(ev.Attachments)
							if err != nil {
								fmt.Println("expected no error unmarshaling attachment with blocks, got: %v", err)
							}
							mess = string(actualAttachmentJson)
						}

						if len(mess) > 0 && mess != "null" && checkID(postId) == true {
							debugLog("User: " + postId + " receive message: " + mess)
							result, ruleInt := checkMessage(mess)

							if reverse == true {
								if result == 0 && ev.Channel != report && ev.Channel != defaultChannel[0] {
									if reacji == true {
										markReaction(api, ev.Channel, ev.TimeStamp, reacjiStr)
									} else {
										postMessageStr(api, defaultChannel[0], defaultChannel[1], mess)
									}
								} else if ev.Channel != report && ev.Channel != defaultChannel[0] {
									markReaction(api, ev.Channel, ev.TimeStamp, label)
								}
							} else {
								if result != 0 && channelMatch(ev.Channel) == false {
									if reacji == true {
										markReaction(api, ev.Channel, ev.TimeStamp, reacjiStr)
									} else {
										postMessage(api, result-1, ruleInt, mess)
									}
									if checkHotline(ruleInt) == true {
										if channelMatch(ev.Channel) == false {
											if reacji == true {
												markReaction(api, ev.Channel, ev.TimeStamp, reacjiStr)
												postMessage(api, result-1, ruleInt, "[Hotline Alert!] "+alertUsers())
											} else {
												postMessage(api, result-1, ruleInt, mess+"\n [Hotline Alert!] "+alertUsers())
											}
										}
									}
								} else if channelMatch(ev.Channel) == false {
									markReaction(api, ev.Channel, ev.TimeStamp, label)
								}
							}
						}

					}
				}
				client.Ack(*evt.Request)
			}
		}
	}()

	go client.Run()
}

func alertUsers() string {
	strs := ""
	for i := 0; i < len(alerts); i++ {
		for r := 0; r < len(alerts[i].USERS); r++ {
			switch alerts[i].USERS[r] {
			case "here":
				strs = strs + " <!here>"
			case "channel":
				strs = strs + " <!channnel>"
			case "everyone":
				strs = strs + " <!everyone>"
			default:
				strs = strs + " <@" + alerts[i].USERS[r] + ">"
			}
		}
	}
	return strs
}

func checkHotline(ruleInt int) bool {
	for i := 0; i < len(alerts); i++ {
		if alerts[i].LABEL == rules[ruleInt].HOTLINE {
			return true
		}
	}

	return false
}

func checkID(ID string) bool {
	for i := 0; i < len(postids); i++ {
		if postids[i] == ID {
			return true
		}
	}
	return false
}

func channelMatch(channel string) bool {
	for i := 0; i < len(incidents); i++ {
		if incidents[i].CHANNNEL == channel {
			return true
		}
	}
	return false
}

func markReaction(api *slack.Client, channnel, ts string, markStr string) {
	msgRef := slack.NewRefToMessage(channnel, ts)

	if err := api.AddReaction(markStr, msgRef); err != nil {
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	}
}

func checkMessage(message string) (int, int) {
	wdays := [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	const layout = "2006/01/02 15:04:05"
	t := time.Now()
	nowDate := t.Format(layout) + " " + wdays[t.Weekday()]

	debugLog("messag: " + message)

	for i := 0; i < len(rules); i++ {
		debugLog("messageRegex: " + rules[i].TARGET)
		messageRegex := regexp.MustCompile(rules[i].TARGET)

		if messageRegex.MatchString(message) == true {
			debugLog("messageRegex: ok")
			debugLog("nowDate: " + nowDate)

			debugLog("dateRegex: " + rules[i].EXCLUDE)
			dateRegex := regexp.MustCompile(rules[i].EXCLUDE)
			if dateRegex.MatchString(nowDate) == true {
				debugLog("dateRegex: ok")
				if act := incidentCheck(rules[i].LABEL); act != 0 {
					return act, i
				}
			}
		}
	}
	return 0, 0
}

func incidentCheck(incidentName string) int {
	for i := 0; i < len(incidents); i++ {
		if incidents[i].LABEL == incidentName {
			return i + 1
		}
	}
	return 0
}
