package slackbot

import (
	"log"

	"github.com/slack-go/slack"
)

func PostMessage(channel, reply, emoji string) {
	_, _, err := ApiBot.PostMessage(channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}
