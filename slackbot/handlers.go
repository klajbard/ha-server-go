package slackbot

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func messageHandler(strArr []string) (string, string) {
	reply := ""
	emoji := ":female-office-worker:"
	switch strArr[0] {
	case "covid":
		infected, dead, cured := getCovidData()
		reply = fmt.Sprintf("*COVID*\n:biohazard_sign: *%d*\n:skull: *%d*\n:heartpulse: *%d*", infected, dead, cured)
		emoji = ":mask:"
	case "help":
		reply = "Type covid to get latest covid data"
	default:
		reply = fmt.Sprintf("Sorry, I dont understand \"_%s_\"", strings.Join(strArr, " "))
	}
	return reply, emoji
}

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
