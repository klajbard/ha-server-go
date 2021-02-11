package main

import (
	"log"
	"net/http"

	"./consumption"
	"./slackbot"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func init() {
	err := godotenv.Load(".ENV")
	if err != nil {
		log.Fatal("Error loading .ENV file")
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
