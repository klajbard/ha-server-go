package slackbot

import (
	"fmt"
	"regexp"
	"strings"

	"../hass"
)

func SetSilence(strArr []string, channel string) {
	re := regexp.MustCompile(`(?i)o(n|ff)`)
	if len(strArr) > 1 && re.Match([]byte(strArr[1])) {
		silence := strArr[1] == "on"

		config := hass.Get()

		config.Silence = !silence
		writeToFile(config)
		PostMessage(channel, fmt.Sprintf("Slack notification is turned %s", strings.ToLower(strArr[1])), ":desktop_computer:")
	}
}
