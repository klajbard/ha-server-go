package slackbot

import (
	"fmt"
	"log"

	"../hass"
	"github.com/slack-go/slack"
)

func handleScraperBlock(value string) {
	switch value {
	case "Bestbuy":
		conf.Enable.Bestbuy = !conf.Enable.Bestbuy
	case "Stockwatcher":
		conf.Enable.Stockwatcher = !conf.Enable.Stockwatcher
	case "Marketplace":
		conf.Enable.Marketplace = !conf.Enable.Marketplace
	case "Steamgifts":
		conf.Enable.Steamgifts = !conf.Enable.Steamgifts
	case "Dht":
		conf.Enable.Dht = !conf.Enable.Dht
	case "Arukereso":
		conf.Enable.Arukereso = !conf.Enable.Arukereso
	case "Covid":
		conf.Enable.Covid = !conf.Enable.Covid
	case "Bumphva":
		conf.Enable.Bumphva = !conf.Enable.Bumphva
	case "Ncore":
		conf.Enable.Ncore = !conf.Enable.Ncore
	case "Fuel":
		conf.Enable.Fuel = !conf.Enable.Fuel
	case "Fixerio":
		conf.Enable.Fixerio = !conf.Enable.Fixerio
	case "Awscost":
		conf.Enable.Awscost = !conf.Enable.Awscost
	case "Btc":
		conf.Enable.Btc = !conf.Enable.Btc
	}

	if conf != nil {
		writeToFile()
	}
}

func sendScraperMessage(channel string) {
	conf = hass.Get()

	btnBestbuy := getButton("Bestbuy", conf.Enable.Bestbuy)
	btnStockwatcher := getButton("Stockwatcher", conf.Enable.Stockwatcher)
	btnMarketplace := getButton("Marketplace", conf.Enable.Marketplace)
	btnSteamgifts := getButton("Steamgifts", conf.Enable.Steamgifts)
	btnDht := getButton("Dht", conf.Enable.Dht)
	btnArukereso := getButton("Arukereso", conf.Enable.Arukereso)
	btnCovid := getButton("Covid", conf.Enable.Covid)
	btnBumphva := getButton("Bumphva", conf.Enable.Bumphva)
	btnNcore := getButton("Ncore", conf.Enable.Ncore)
	btnFuel := getButton("Fuel", conf.Enable.Fuel)
	btnFixerio := getButton("Fixerio", conf.Enable.Fixerio)
	btnAwscost := getButton("Awscost", conf.Enable.Awscost)
	btnBtc := getButton("Btc", conf.Enable.Btc)
	actionBlock := slack.NewActionBlock("scraper", btnBestbuy, btnStockwatcher, btnMarketplace, btnSteamgifts, btnDht, btnArukereso, btnCovid, btnBumphva, btnFuel, btnNcore, btnFixerio, btnAwscost, btnBtc)

	btnText := slack.NewTextBlockObject("plain_text", "~ Done ~", false, false)
	btn := slack.NewButtonBlockElement("", "done", btnText)
	btnBlock := slack.NewActionBlock("hassio", btn)

	_, _, err := ApiBot.PostMessage(channel, slack.MsgOptionBlocks(actionBlock, btnBlock), slack.MsgOptionIconEmoji(":construction_worker:"))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}

func getButton(text string, value bool) *slack.ButtonBlockElement {
	displayText := fmt.Sprintf("%s: %t", text, value)
	btnText := slack.NewTextBlockObject("plain_text", displayText, false, false)
	btn := slack.NewButtonBlockElement("", text, btnText)
	return btn
}
