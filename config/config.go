package config

import (
	"os"

	"github.com/slack-go/slack"
)

var ApiBot *slack.Client
var ApiUser *slack.Client

func init() {

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	userToken := os.Getenv("SLACK_OAUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	ApiBot = slack.New(botToken, slack.OptionAppLevelToken(appToken))
	ApiUser = slack.New(userToken, slack.OptionAppLevelToken(appToken))
}
