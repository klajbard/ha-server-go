package slackbot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var ApiBot *slack.Client
var ApiUser *slack.Client

func Run() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	userToken := os.Getenv("SLACK_OAUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	ApiBot = slack.New(botToken, slack.OptionAppLevelToken(appToken))
	ApiUser = slack.New(userToken, slack.OptionAppLevelToken(appToken))

	client := socketmode.New(ApiBot)

	go handleEvents(client)

	client.Run()
}

func callbackMux(callback slack.InteractionCallback) {
	timestamp := callback.Container.MessageTs
	channel := callback.Channel.GroupConversation.Conversation.ID
	value := callback.ActionCallback.BlockActions[0].Value

	if value == "hum" {
		Humidity(channel)
	} else if value == "temp" {
		Temperature(channel)
	} else if value == "covid" {
		Covid(channel)
	}
	_, _, err := ApiUser.DeleteMessage(channel, timestamp)
	if err != nil {
		log.Printf("Deleting message failed: %v", err)
	}
}

func handleEvents(client *socketmode.Client) {
	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeInteractive:
			callback, ok := evt.Data.(slack.InteractionCallback)
			if !ok {
				log.Printf("Something bad hapened: %+v\n", evt)
				continue
			}

			if callback.Type == slack.InteractionTypeBlockActions {
				callbackMux(callback)
			}
			client.Ack(*evt.Request)

		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				continue
			}
			client.Ack(*evt.Request)

			if eventsAPIEvent.Type == slackevents.CallbackEvent {
				eventMux(eventsAPIEvent)
			}
		}
	}
}

func eventMux(eventsAPIEvent slackevents.EventsAPIEvent) {
	switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
	case *slackevents.MessageEvent:
		if strArr := strings.Split(ev.Text, " "); len(strArr) > 1 {
			re := regexp.MustCompile(`(?i)athena`) // Bot name
			match := re.Match([]byte(strArr[0]))
			if strArr[0] == fmt.Sprintf("<@%s>", os.Getenv("SLACK_BOT_ID")) || match {
				messageArr := strArr[1:]
				messageMux(messageArr, ev.Channel)
				_, _, err := ApiUser.DeleteMessage(ev.Channel, ev.TimeStamp)
				if err != nil {
					log.Printf("Deleting message failed: %v", err)
				}
			}
		}
	}
}

func messageMux(strArr []string, channel string) {
	switch strArr[0] {
	case "cons":
		Consumption(strArr, channel)
	case "covid":
		Covid(channel)
	case "hum":
		Humidity(channel)
	case "temp":
		Temperature(channel)
	case "turn":
		TurnSwitch(strArr, channel)
	case "arukereso":
		Arukereso(strArr, channel)
	case "hautils":
		IsRunning(channel)
	case "help":
		Help(channel)
	case "commands":
		sendBlockMessages(channel)
	default:
		Default(strArr, channel)
	}
}

func sendBlockMessages(channel string) {
	tempBtnText := slack.NewTextBlockObject("plain_text", "Temperature", false, false)
	humBtnText := slack.NewTextBlockObject("plain_text", "Humidity", false, false)
	covidBtnText := slack.NewTextBlockObject("plain_text", "Covid", false, false)
	tempBtn := slack.NewButtonBlockElement("", "temp", tempBtnText)
	humBtn := slack.NewButtonBlockElement("", "hum", humBtnText)
	covidBtn := slack.NewButtonBlockElement("", "covid", covidBtnText)
	actionBlock := slack.NewActionBlock("", tempBtn, humBtn, covidBtn)

	_, _, err := ApiBot.PostMessage(channel, slack.MsgOptionBlocks(actionBlock))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}
