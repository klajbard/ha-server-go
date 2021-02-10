package main

import (
	"net/http"

	"./consumption"
	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()
	mux.GET("/cons/:device", consumption.GetAllCons)
	mux.GET("/cons/:device/:date", consumption.GetCons)
	mux.PUT("/cons/:device", consumption.PutCons)
	http.ListenAndServe(":5500", mux)
}
