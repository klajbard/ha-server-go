package slackbot

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../config"
	"../types"
	"github.com/PuerkitoBio/goquery"
)

func Arukereso(strArr []string, channel string) {
	var sb strings.Builder

	reply := "Wrong parameters"
	emoji := ":desktop_computer:"

	if len(strArr) < 2 {
		PostMessage(channel, reply, emoji)
	}

	for _, item := range config.Conf.Arukereso {
		if item.Name == strArr[1] {
			for _, url := range item.Urls {
				result := queryAK(url)
				sb.WriteString(fmt.Sprintf("<%s|*%s - %d*>\n", url, result.Name, result.Price))
			}
		}
	}

	if sb.Len() > 0 {
		reply = sb.String()
		PostMessage(channel, reply, emoji)
	}
}

func queryAK(url string) types.ArukeresoResult {
	item := types.ArukeresoResult{}
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return item
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return item
	}

	item.Name = strings.TrimSpace(doc.Find("h1.hidden-xs").Text())

	doc.Find(".optoffer.device-desktop").Each(func(_ int, s *goquery.Selection) {
		price, _ := s.Find("[itemprop=\"price\"]").Attr("content")
		priceInt, _ := strconv.Atoi(price)
		if item.Price == 0 || priceInt < item.Price {
			item.Price = priceInt
		}
	})
	return item
}
