package slackbot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"../config"
	"../handlers"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func Run() {
	client := socketmode.New(config.ApiBot)

	go handleEvents(client)

	client.Run()
}

func callbackMux(callback slack.InteractionCallback) {
	timestamp := callback.Container.MessageTs
	channel := callback.Channel.GroupConversation.Conversation.ID
	value := callback.ActionCallback.BlockActions[0].Value

	if value == "hum" {
		handlers.Humidity(channel)
	} else if value == "temp" {
		handlers.Temperature(channel)
	} else if value == "covid" {
		handlers.Covid(channel)
	}
	_, _, err := config.ApiUser.DeleteMessage(channel, timestamp)
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
				_, _, err := config.ApiUser.DeleteMessage(ev.Channel, ev.TimeStamp)
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
		handlers.Consumption(strArr, channel)
	case "covid":
		handlers.Covid(channel)
	case "hum":
		handlers.Humidity(channel)
	case "temp":
		handlers.Temperature(channel)
	case "turn":
		handlers.TurnSwitch(strArr, channel)
	case "arukereso":
		handlers.Arukereso(strArr, channel)
	case "hautils":
		handlers.IsRunning(channel)
	case "help":
		handlers.Help(channel)
	case "commands":
		sendBlockMessages(channel)
	default:
		handlers.Default(strArr, channel)
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

	_, _, err := config.ApiBot.PostMessage(channel, slack.MsgOptionBlocks(actionBlock))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}
