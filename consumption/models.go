package consumption

import (
	"github.com/klajbard/ha-server-go/config"
	"gopkg.in/mgo.v2/bson"
)

type Consumption struct {
	Device string  // `json:"device" bson:"device"`
	Date   string  // `json:"date" bson:"date"`
	Watt   float64 // `json:"watt" bson:"watt"`
}

func AllCons(d string) ([]Consumption, error) {
	cons := []Consumption{}

	err := config.Consumptions.Find(bson.M{"device": d}).All(&cons)
	if err != nil {
		return nil, err
	}
	return cons, nil
}

func OneCons(d string, date string) Consumption {
	cons := Consumption{}

	err := config.Consumptions.Find(bson.M{"device": d, "date": date}).One(&cons)
	if err != nil {
		return Consumption{
			Device: d,
			Date:   date,
			Watt:   0.0,
		}
	}
	return cons
}
