package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"../config"
	"../hass"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"gopkg.in/yaml.v2"
)

var ApiBot *slack.Client
var ApiUser *slack.Client
var conf *hass.Configuration

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

func writeToFile() {
	output, err := yaml.Marshal(conf)

	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile(config.Conf.ScraperConfig, output, 0)

	if err != nil {
		log.Println(err)
	}
}

func removeMessage(channel, timestamp string) {
	_, _, err := ApiUser.DeleteMessage(channel, timestamp)
	if err != nil {
		log.Printf("Deleting message failed: %v", err)
	}
}

func callbackMux(callback slack.InteractionCallback) {
	timestamp := callback.Container.MessageTs
	channel := callback.Channel.GroupConversation.Conversation.ID
	value := callback.ActionCallback.BlockActions[0].Value
	block := callback.ActionCallback.BlockActions[0].BlockID

	switch block {
	case "scraper":
		handleScraperBlock(value)
		removeMessage(channel, timestamp)
		if value != "done" {
			sendScraperMessage(channel)
		}
	case "hassio":
		handleHassioBlock(value, channel)
		removeMessage(channel, timestamp)
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
		sendHassioMessage(channel)
	case "scraper":
		sendScraperMessage(channel)
	default:
		Default(strArr, channel)
	}
}
