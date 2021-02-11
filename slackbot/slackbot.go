package slackbot

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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
							_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText(strings.Join(strArr[1:], " "), false))
							if err != nil {
								fmt.Printf("failed posting message: %v", err)
							}
						}
					}
				}
			}
		}
	}
}
