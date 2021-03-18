package slackbot

import (
	"fmt"
	"os/exec"
	"strings"
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

func Default(strArr []string, channel string) {
	reply := fmt.Sprintf("Sorry, I dont understand \"_%s_\"", strings.Join(strArr, " "))
	emoji := ":female-office-worker:"

	PostMessage(channel, reply, emoji)
}
