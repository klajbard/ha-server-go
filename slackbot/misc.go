package slackbot

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/slack-go/slack"
)

func IsRunning() bool {
	_, err := exec.Command("/bin/systemctl", "is-active", "--quiet", "hautils.service").Output()
	return err == nil
}

func StartService(strArr []string, channel string) {
	reply := "Wrong parameters"
	emoji := ":female-office-worker:"
	if len(strArr) >= 2 {
		service := strArr[1]
		reply = fmt.Sprintf("%s is started successfully", service)
		command := fmt.Sprintf("sudo systemctl start %s", service)
		_, err := exec.Command("/bin/sh", "-c", command).Output()
		if err != nil {
			reply = fmt.Sprintf("Something bad happened: %s", err.Error())
		}
	}
	PostMessage(channel, reply, emoji)
}

func StopService(strArr []string, channel string) {
	reply := "Wrong parameters"
	emoji := ":female-office-worker:"
	if len(strArr) >= 2 {
		service := strArr[1]
		reply = fmt.Sprintf("%s is stopped successfully", service)
		command := fmt.Sprintf("sudo systemctl stop %s", service)
		_, err := exec.Command("/bin/sh", "-c", command).Output()
		if err != nil {
			reply = fmt.Sprintf("Something bad happened: %s", err.Error())
		}
	}
	PostMessage(channel, reply, emoji)
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
*aklist <add/rm> <url>* - add/remove url to arukereso query list
*arukereso <product>* - show current lowest price for _product_
*cons <sensor>* - sensor consumption
*covid* - current covid data
*hassio* - show hassio sensors
*hautils* - check if scraper(hautils) is running
*notif <on/off>* - sets whether hautils should notify via slack
*hum* - display humidity
*start/stop <service> - start/stop systemd service*
*scraper* - set and display current scrapers status
*temp* - display temperature
*turn <sensor> <on/off>* - turn switch on/off`
	emoji := ":female-office-worker:"

	PostMessage(channel, reply, emoji)
}

func SendEmpty(channel string) {
	statusText := slack.NewTextBlockObject("plain_text", "Available commands", false, false)
	headerSection := slack.NewHeaderBlock(statusText, slack.HeaderBlockOptionBlockID("test_block"))

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
