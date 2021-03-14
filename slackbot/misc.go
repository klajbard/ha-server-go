package slackbot

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsRunning(channel string) {
	reply := "Ha utils is running"
	emoji := ":female-office-worker:"
	_, err := exec.Command("/bin/systemctl", "is-active", "--quiet", "hautils.service").Output()
	if err != nil {
		reply = "Ha utils is not running"
	}

	PostMessage(channel, reply, emoji)
}

func Help(channel string) {
	reply := `
*covid* - current covid data
*cons <sensor>* - sensor consumption
*hum* - display humidity
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
