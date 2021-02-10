package config

import (
	"log"

	"gopkg.in/mgo.v2"
)

var DB *mgo.Database
var Consumptions *mgo.Collection

func init() {
	s, err := mgo.Dial("mongodb://localhost:27017/hassio")
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("hassio")
	Consumptions = DB.C("consumptions")

	log.Println("MongoDB connected")
}
