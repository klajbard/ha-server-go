package slackbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"../consumption"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SensorData struct {
	EntityId   string   // `json:"entity_id" bson:"entity_id"`
	State      string   // `json:"state" bson:"state"`
	Attributes struct { // `json:"attributes" bson:"attributes"`
		UnitOfMeasurement string // `json:"unit_of_measurement" bson:"unit_of_measurement"`
		FriendlyName      string // `json:"friendly_name" bson:"friendly_name"`
		DeviceClass       string // `json:"device_class" bson:"device_class"`
	}
	LastChanged string   // `json:"last_changed" bson:"last_changed"`
	LastUpdated string   // `json:"last_updated" bson:"last_updated"`
	Context     struct { // `json:"context" bson:"context"`
		Id       string // `json:"id" bson:"id"`
		ParentId string // `json:"device_class" bson:"device_class"`
		UserId   string // `json:"user_id" bson:"user_id"`
	}
}

type SensorValue struct {
	Name  string
	Value string
}

var HUMIDITIES = map[string]string{
	"sensor.xiaomi_airpurifier_humidity": "Purifier humidity",
	"sensor.xiaomi_humidifier_humidity":  "Humidifier humidity",
	"sensor.rpi_humidity":                "Raspberry humidity",
	"sensor.aqara_temp_humidity":         "Aqara bedroom humidity",
	"sensor.aqara_temp2_humidity":        "Aqara living room humidity",
	"sensor.aqara_temp3_humidity":        "Aqara balcony humidity",
	"sensor.mijia_temp_humidity":         "Mijia humidity",
}

var TEMPERATURES = map[string]string{
	"sensor.xiaomi_airpurifier_temp": "Purifier temperature",
	"sensor.xiaomi_humidifier_temp":  "Humidifier temperature",
	"sensor.rpi_temperature":         "Raspberry temperature",
	"sensor.aqara_temp_temperature":  "Aqara bedroom temperature",
	"sensor.aqara_temp2_temperature": "Aqara living room temperature",
	"sensor.aqara_temp3_temperature": "Aqara balcony temperature",
	"sensor.mijia_temp_temperature":  "Mijia temperature",
	"sensor.mandula_temp":            "Mandula temperature",
}

func Run() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	api := slack.New(botToken, slack.OptionAppLevelToken(appToken))
	client := socketmode.New(api)

	go handleEvents(api, client)

	client.Run()
}

func handleEvents(api *slack.Client, client *socketmode.Client) {
	for evt := range client.Events {
		if evt.Type == socketmode.EventTypeEventsAPI {
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				continue
			}
			client.Ack(*evt.Request)

			if eventsAPIEvent.Type == slackevents.CallbackEvent {
				eventMux(api, eventsAPIEvent)
			}
		}
	}
}

func eventMux(api *slack.Client, eventsAPIEvent slackevents.EventsAPIEvent) {
	switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
	case *slackevents.MessageEvent:
		if strArr := strings.Split(ev.Text, " "); len(strArr) > 1 {
			re := regexp.MustCompile(`(?i)athena`) // Bot name
			match := re.Match([]byte(strArr[0]))
			if strArr[0] == fmt.Sprintf("<@%s>", os.Getenv("SLACK_BOT_ID")) || match {
				messageArr := strArr[1:]
				reply, emoji := messageMux(messageArr)
				_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
				if err != nil {
					log.Printf("Posting message failed: %v", err)
				}
			}
		}
	}
}

func messageMux(strArr []string) (string, string) {
	reply := ""
	emoji := ":female-office-worker:"
	switch strArr[0] {
	case "covid":
		infected, dead, cured := getCovidData()
		reply = fmt.Sprintf("*COVID*\n:biohazard_sign: *%d*\n:skull: *%d*\n:heartpulse: *%d*", infected, dead, cured)
		emoji = ":mask:"
	case "cons":
		if len(strArr) < 2 {
			reply = "Please specify a sensor name: *cons <sensor>*"
		} else {
			today := time.Now().Format("20160102")
			cons := consumption.OneCons(strArr[1], today)
			reply = fmt.Sprintf("*%s* today's consumption: *%.2f Wh*", cons.Device, cons.Watt)
		}
	case "hum":
		if len(strArr) < 2 || strArr[1] == "all" {
			var sb strings.Builder
			hums := GetAllHumidity()
			for _, sensor := range hums {
				sb.WriteString(fmt.Sprintf("*%s*: %s %%\n", sensor.Name, sensor.Value))
			}
			emoji = ":droplet:"
			reply = sb.String()
		}
	case "temp":
		if len(strArr) < 2 || strArr[1] == "all" {
			var sb strings.Builder
			hums := GetAllTemp()
			for _, sensor := range hums {
				sb.WriteString(fmt.Sprintf("*%s*: %s °C\n", sensor.Name, sensor.Value))
			}
			emoji = ":thermometer:"
			reply = sb.String()
		}
	case "help":
		reply = `
*covid*					current covid data
*cons <sensor>*	sensor consumption
*hum*						display humidity
*temp*					display temperature`
	default:
		reply = fmt.Sprintf("Sorry, I dont understand \"_%s_\"", strings.Join(strArr, " "))
	}
	return reply, emoji
}

func GetAllTemp() []SensorValue {
	ret := []SensorValue{}
	for sensor, name := range TEMPERATURES {
		ret = append(ret, SensorValue{name, getHassioData(sensor)})
	}

	return ret
}

func GetAllHumidity() []SensorValue {
	ret := []SensorValue{}
	for sensor, name := range HUMIDITIES {
		ret = append(ret, SensorValue{name, getHassioData(sensor)})
	}

	return ret
}

func getHassioData(sensor string) string {
	sensorData := &SensorData{}
	req, err := http.NewRequest("GET", "http://192.168.1.27:8123/api/states/"+sensor, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("HASS_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal([]byte(string(body)), sensorData)
	if err != nil {
		log.Println(err)
	}

	return sensorData.State
}
