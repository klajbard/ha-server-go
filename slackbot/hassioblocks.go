package slackbot

import (
	"log"

	"github.com/slack-go/slack"
)

func handleHassioBlock(value, channel string) {
	if value == "hum" {
		Humidity(channel)
	} else if value == "temp" {
		Temperature(channel)
	} else if value == "covid" {
		Covid(channel)
	}
}

func sendHassioMessage(channel string) {
	tempBtnText := slack.NewTextBlockObject("plain_text", "Temperature", false, false)
	humBtnText := slack.NewTextBlockObject("plain_text", "Humidity", false, false)
	covidBtnText := slack.NewTextBlockObject("plain_text", "Covid", false, false)
	tempBtn := slack.NewButtonBlockElement("", "temp", tempBtnText)
	humBtn := slack.NewButtonBlockElement("", "hum", humBtnText)
	covidBtn := slack.NewButtonBlockElement("", "covid", covidBtnText)
	actionBlock := slack.NewActionBlock("hassio", tempBtn, humBtn, covidBtn)

	_, _, err := ApiBot.PostMessage(channel, slack.MsgOptionBlocks(actionBlock))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}
