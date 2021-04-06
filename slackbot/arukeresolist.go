package slackbot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/klajbard/ha-server-go/hass"
)

func isValidUrl(url string) bool {
	re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	return re.Match([]byte(url))
}

func addItem(item, channel string) {
	item = strings.Split(strings.Trim(item, "<"), "|")[0]
	if isValidUrl(item) {
		config := hass.Get()
		url := hass.Url{
			Url: item,
		}
		config.Arukereso = append(config.Arukereso, url)
		writeToFile(config)
		PostMessage(channel, fmt.Sprintf("%s is added to the list", item), ":desktop_computer:")
	} else {
		PostMessage(channel, fmt.Sprintf("%s is not a valid url", item), ":desktop_computer:")
	}
}

func removeItem(item, channel string) {
	item = strings.Split(strings.Trim(item, "<>"), "|")[0]
	config := hass.Get()
	cleaned := make([]hass.Url, len(config.Arukereso))
	idx := 0
	for _, url := range config.Arukereso {
		if item != url.Url {
			cleaned[idx] = url
			idx++
		}
	}

	config.Arukereso = cleaned[:idx]
	writeToFile(config)
	PostMessage(channel, fmt.Sprintf("%s is removed from the list", item), ":desktop_computer:")
}

func AKList(strArr []string, channel string) {
	if len(strArr) > 2 {
		switch strArr[1] {
		case "add":
			addItem(strArr[2], channel)
			return
		case "rm":
			removeItem(strArr[2], channel)
			return
		}
	}
	var sb strings.Builder
	conf := hass.Get()

	for _, item := range conf.Arukereso {
		sb.WriteString(fmt.Sprintf("- %s\n", item.Url))
	}
	PostMessage(channel, sb.String(), ":desktop_computer:")
}
