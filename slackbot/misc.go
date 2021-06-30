package slackbot

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

func IsRunning(service string) bool {
	_, err := exec.Command("/bin/systemctl", "is-active", "--quiet", fmt.Sprintf("%s.service", service)).Output()
	return err == nil
}

func AkGoQueryStatus() (string, error) {
	output, err := exec.Command("/usr/bin/journalctl", "-au", "akgoquery.service", "-n", "1").Output()
	if err != nil {
		return "", err
	}
	outStr := string(output)
	if strings.Contains(outStr, "Succeeded.") {
		return "Query job is successfully done.", nil
	} else if strings.Contains(outStr, "Querying...") {
		r := regexp.MustCompile(`\d+/\d+`)
		matches := r.FindAllString(outStr, -1)
		return fmt.Sprintf("Query job is running %s", matches[0]), nil
	} else {
		return "Query job failed", nil
	}
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

func SendAkGoQueryStatus(channel string) {
	output, err := AkGoQueryStatus()
	emoji := ":female-office-worker:"

	if err != nil {
		fmt.Println(err)
	} else {
		PostMessage(channel, output, emoji)
	}
}

func StatusService(strArr []string, channel string) {
	reply := "Wrong parameters"
	emoji := ":female-office-worker:"
	if len(strArr) >= 2 {
		service := strArr[1]
		reply = fmt.Sprintf("%s is running", service)
		if !IsRunning(service) {
			reply = fmt.Sprintf("%s is not running", service)
		}
	}

	PostMessage(channel, reply, emoji)
}

func Help(channel string) {
	reply := `
*akgostatus* - Status of the akgoquery systemd status
*aklist <add/rm> <url>* - add/remove url to arukereso query list
*arukereso <product>* - show current lowest price for _product_
*cons <sensor>* - sensor consumption
*covid* - current covid data
*hassio* - show hassio sensors
*notif <on/off>* - sets whether hautils should notify via slack
*hum* - display humidity
*status <service>* - check if systemd is running
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
