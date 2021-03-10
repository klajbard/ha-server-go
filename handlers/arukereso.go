package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../types"
	"../utils"
	"github.com/PuerkitoBio/goquery"
)

var ARUKERESO_URLS_CPU = []string{
	"https://www.arukereso.hu/processzor-c3139/intel/core-i5-10400f-6-core-2-9ghz-lga1200-p558582354/",
	"https://www.arukereso.hu/processzor-c3139/intel/core-i5-10400-6-core-2-9ghz-lga1200-p558582279/",
	"https://www.arukereso.hu/processzor-c3139/intel/core-i5-10500-6-core-3-1ghz-lga1200-p558586827/",
	"https://www.arukereso.hu/processzor-c3139/intel/core-i5-10600k-6-core-4-1ghz-lga1200-p558587868/",
}

var ARUKERESO_URLS_ALAPLAP = []string{
	"https://www.arukereso.hu/alaplap-c3128/asrock/b560m-pro4-p633866961",
}

func Arukereso(strArr []string, channel string) {
	var sb strings.Builder

	reply := "Wrong parameters"
	emoji := ":desktop_computer:"

	if len(strArr) < 2 {
		utils.PostMessage(channel, reply, emoji)
	}

	switch strArr[1] {
	case "proci":
		for _, url := range ARUKERESO_URLS_CPU {
			result := queryAK(url)
			sb.WriteString(fmt.Sprintf("<%s|*%s - %d*>\n", url, result.Name, result.Price))
		}
	case "alaplap":
		for _, url := range ARUKERESO_URLS_ALAPLAP {
			result := queryAK(url)
			sb.WriteString(fmt.Sprintf("<%s|*%s - %d*>\n", url, result.Name, result.Price))
		}
	}
	reply = sb.String()

	utils.PostMessage(channel, reply, emoji)
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
