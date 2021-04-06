package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/klajbard/ha-server-go/consumption"
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
	go slackbot.Run()

	mux := httprouter.New()
	mux.GET("/cons/:device", consumption.GetAllCons)
	mux.GET("/cons/:device/:date", consumption.GetCons)
	mux.PUT("/cons/:device", consumption.PutCons)
	http.ListenAndServe(":5500", mux)
}
