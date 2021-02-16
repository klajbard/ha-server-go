package slackbot

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getCovidData() (int, int, int) {
	resp, err := http.Get("https://koronavirus.gov.hu")
	if err != nil {
		log.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	infectedPest := getNum(doc.Find("#api-fertozott-pest").Text())
	infectedVidek := getNum(doc.Find("#api-fertozott-videk").Text())
	deadPest := getNum(doc.Find("#api-elhunyt-pest").Text())
	deadVidek := getNum(doc.Find("#api-elhunyt-videk").Text())
	curedPest := getNum(doc.Find("#api-gyogyult-pest").Text())
	curedVidek := getNum(doc.Find("#api-gyogyult-videk").Text())

	infected := infectedPest + infectedVidek
	dead := deadPest + deadVidek
	cured := curedPest + curedVidek
	return infected, dead, cured
}

func getNum(input string) int {
	trimmed := strings.ReplaceAll(input, " ", "")
	num, err := strconv.Atoi(trimmed)
	if err != nil {
		log.Println(err)
	}

	return num
}
