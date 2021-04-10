package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/klajbard/ha-server-go/slackbot"
)

func init() {
	err := godotenv.Load(".ENV")
	if err != nil {
		log.Println("Error loading .ENV file")
	}
	log.Println("Loaded env variables")
}

func main() {
	slackbot.Run()
}
