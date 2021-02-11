package slackbot

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func Run() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	api := slack.New(botToken, slack.OptionAppLevelToken(appToken))
	client := socketmode.New(api)

	go handleEvents(api, client)

	client.Run()
}

func handleEvents(api *slack.Client, client *socketmode.Client) {
	for evt := range client.Events {
		if evt.Type == socketmode.EventTypeEventsAPI {
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				continue
			}
			client.Ack(*evt.Request)

			if eventsAPIEvent.Type == slackevents.CallbackEvent {
				switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
				case *slackevents.MessageEvent:
					if strArr := strings.Split(ev.Text, " "); len(strArr) > 1 {
						re := regexp.MustCompile(`(?i)hestia`)
						match := re.Match([]byte(strArr[0]))
						if strArr[0] == fmt.Sprintf("<@%s>", os.Getenv("SLACK_BOT_ID")) || match {
							messageArr := strArr[1:]
							reply, emoji := messageHandler(messageArr)
							_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
							if err != nil {
								log.Printf("Posting message failed: %v", err)
							}
						}
					}
				}
			}
		}
	}
}

func messageHandler(strArr []string) (string, string) {
	reply := ""
	emoji := ":female-office-worker:"
	switch strArr[0] {
	case "covid":
		infected, dead, cured := getCovidData()
		reply = fmt.Sprintf("*COVID*\n:biohazard_sign: *%d*\n:skull: *%d*\n:heartpulse: *%d*", infected, dead, cured)
		emoji = ":mask:"
	case "help":
		reply = "Type covid to get latest covid data"
	default:
		reply = fmt.Sprintf("Sorry, I dont understand \"_%s_\"", strings.Join(strArr, " "))
	}
	return reply, emoji
}

func getCovidData() (int, int, int) {
	resp, err := http.Get("https://koronavirus.gov.hu")
	if err != nil {
		log.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	infectedPest := getNum(doc.Find("#api-fertozott-pest").Text())
	infectedVidek := getNum(doc.Find("#api-fertozott-videk").Text())
	deadPest := getNum(doc.Find("#api-elhunyt-pest").Text())
	deadVidek := getNum(doc.Find("#api-elhunyt-videk").Text())
	curedPest := getNum(doc.Find("#api-gyogyult-pest").Text())
	curedVidek := getNum(doc.Find("#api-gyogyult-videk").Text())

	infected := infectedPest + infectedVidek
	dead := deadPest + deadVidek
	cured := curedPest + curedVidek
	return infected, dead, cured
}

func getNum(input string) int {
	trimmed := strings.ReplaceAll(input, " ", "")
	num, err := strconv.Atoi(trimmed)
	if err != nil {
		log.Println(err)
	}

	return num
}
