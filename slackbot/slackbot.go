package slackbot

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"../config"
	"../consumption"
	"../handlers"
	"../utils"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func Run() {
	client := socketmode.New(config.ApiBot)

	go handleEvents(client)

	client.Run()
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

			timestamp := callback.Container.MessageTs
			channel := callback.Channel.GroupConversation.Conversation.ID
			value := callback.ActionCallback.BlockActions[0].Value

			var payload interface{}

			if callback.Type == slack.InteractionTypeBlockActions {
				if value == "hum" {
					sendHumidity(channel)
				} else if value == "temp" {
					sendTemperature(channel)
				} else if value == "covid" {
					sendCovid(channel)
				}
				_, _, err := config.ApiUser.DeleteMessage(channel, timestamp)
				if err != nil {
					log.Printf("Deleting message failed: %v", err)
				}
			}
			client.Ack(*evt.Request, payload)

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
				if messageArr[0] == "commands" {
					sendBlockMessages(ev.Channel)
				} else {
					messageMux(messageArr, ev.Channel)
				}
				_, _, err := config.ApiUser.DeleteMessage(ev.Channel, ev.TimeStamp)
				if err != nil {
					log.Printf("Deleting message failed: %v", err)
				}
			}
		}
	}
}

func messageMux(strArr []string, channel string) {
	reply := ""
	emoji := ":female-office-worker:"
	switch strArr[0] {
	case "cons":
		if len(strArr) < 2 {
			reply = "Please specify a sensor name: *cons <sensor>*"
		} else {
			today := time.Now().Format("06.01.02")
			cons := consumption.OneCons(strArr[1], today)
			reply = fmt.Sprintf("*%s* today's consumption: *%.2f Wh*", cons.Device, cons.Watt)
		}
	case "covid":
		sendCovid(channel)
	case "hum":
		sendHumidity(channel)
	case "temp":
		sendTemperature(channel)
	case "turn":
		handlers.TurnSwitch(strArr, channel)
	case "arukereso":
		handlers.Arukereso(strArr, channel)
	case "hautils":
		_, err := exec.Command("/bin/systemctl", "is-active", "--quiet", "hautils.service").Output()
		if err != nil {
			reply = "Ha utils is not running"
		} else {
			reply = "Ha utils is running"
		}

	case "help":
		reply = `
*covid* - current covid data
*cons <sensor>* - sensor consumption
*hum* - display humidity
*temp* - display temperature
*turn <sensor> <on/off>* - turn switch on/off`
	default:
		reply = fmt.Sprintf("Sorry, I dont understand \"_%s_\"", strings.Join(strArr, " "))
	}

	_, _, err := config.ApiBot.PostMessage(channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
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

func sendHumidity(channel string) {
	var sb strings.Builder
	hums := handlers.GetAllHumidity()
	for _, sensor := range hums {
		sb.WriteString(fmt.Sprintf("*%s*: %s %%\n", sensor.Name, sensor.Value))
	}
	emoji := ":droplet:"
	reply := sb.String()

	utils.PostMessage(channel, reply, emoji)
}

func sendTemperature(channel string) {
	var sb strings.Builder
	hums := handlers.GetAllTemp()
	for _, sensor := range hums {
		sb.WriteString(fmt.Sprintf("*%s*: %s °C\n", sensor.Name, sensor.Value))
	}
	emoji := ":thermometer:"
	reply := sb.String()

	utils.PostMessage(channel, reply, emoji)
}

func sendCovid(channel string) {
	infected, dead, cured := getCovidData()
	reply := fmt.Sprintf("*COVID*\n:biohazard_sign: *%d*\n:skull: *%d*\n:heartpulse: *%d*", infected, dead, cured)
	emoji := ":mask:"

	utils.PostMessage(channel, reply, emoji)
}
