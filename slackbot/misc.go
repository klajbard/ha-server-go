package slackbot

import (
	"log"
	"os/exec"

	"github.com/slack-go/slack"
)

func IsRunning() bool {
	_, err := exec.Command("/bin/systemctl", "is-active", "--quiet", "hautils.service").Output()
	return err == nil
}

func SendIsRunning(channel string) {
	reply := "Ha utils is running"
	emoji := ":female-office-worker:"
	if !IsRunning() {
		reply = "Ha utils is not running"
	}

	PostMessage(channel, reply, emoji)
}

func Help(channel string) {
	reply := `
*arukereso <product>* - show current lowest price for _product_
*cons <sensor>* - sensor consumption
*covid* - current covid data
*hassio* - show hassio sensors
*hautils* - check if scraper(hautils) is running
*hum* - display humidity
*scraper* - set and display current scrapers status
*temp* - display temperature
*turn <sensor> <on/off>* - turn switch on/off`
	emoji := ":female-office-worker:"

	PostMessage(channel, reply, emoji)
}

func SendEmpty(channel string) {
	statusText := slack.NewTextBlockObject("plain_text", "Available commands", false, false)
	headerSection := slack.NewSectionBlock(statusText, nil, nil)

	btncovid := getSimpleButton("covid")
	btnhassio := getSimpleButton("hassio")
	btnhautils := getSimpleButton("hautils")
	btnhum := getSimpleButton("hum")
	btnscraper := getSimpleButton("scraper")
	btntemp := getSimpleButton("temp")
	btncancel := getSimpleButton("cancel")
	actionBlock := slack.NewActionBlock("commands", btncancel, btncovid, btnhassio, btnhautils, btnhum, btnscraper, btntemp)

	_, _, err := ApiBot.PostMessage(channel, slack.MsgOptionBlocks(headerSection, actionBlock), slack.MsgOptionIconEmoji(":female-office-worker:"))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}

func getSimpleButton(text string) *slack.ButtonBlockElement {
	btnText := slack.NewTextBlockObject("plain_text", text, false, false)
	btn := slack.NewButtonBlockElement("", text, btnText)
	return btn
}
