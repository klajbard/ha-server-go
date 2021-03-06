package slackbot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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

var apiBot *slack.Client
var apiUser *slack.Client

func Run() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	userToken := os.Getenv("SLACK_OAUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	apiBot = slack.New(botToken, slack.OptionAppLevelToken(appToken))
	apiUser = slack.New(userToken, slack.OptionAppLevelToken(appToken))

	client := socketmode.New(apiBot)

	go handleEvents(client)

	client.Run()
}

func handleEvents(client *socketmode.Client) {
	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeInteractive:
			callback, ok := evt.Data.(slack.InteractionCallback)
			if !ok {
				log.Printf("Something bad hapened: %+v\n", evt)
				continue
			}

			timestamp := callback.Container.MessageTs
			channel := callback.Channel.GroupConversation.Conversation.ID
			value := callback.ActionCallback.BlockActions[0].Value

			var payload interface{}

			if callback.Type == slack.InteractionTypeBlockActions {
				if value == "hum" {
					sendHumidity(channel)
				} else if value == "temp" {
					sendTemperature(channel)
				} else if value == "covid" {
					sendCovid(channel)
				}
				_, _, err := apiUser.DeleteMessage(channel, timestamp)
				if err != nil {
					log.Printf("Deleting message failed: %v", err)
				}
			}
			client.Ack(*evt.Request, payload)

		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				continue
			}
			client.Ack(*evt.Request)

			if eventsAPIEvent.Type == slackevents.CallbackEvent {
				eventMux(eventsAPIEvent)
			}
		}
	}
}

func eventMux(eventsAPIEvent slackevents.EventsAPIEvent) {
	switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
	case *slackevents.MessageEvent:
		if strArr := strings.Split(ev.Text, " "); len(strArr) > 1 {
			re := regexp.MustCompile(`(?i)athena`) // Bot name
			match := re.Match([]byte(strArr[0]))
			if strArr[0] == fmt.Sprintf("<@%s>", os.Getenv("SLACK_BOT_ID")) || match {
				messageArr := strArr[1:]
				if messageArr[0] == "commands" {
					sendBlockMessages(ev.Channel)
				} else {
					messageMux(messageArr, ev.Channel)
				}
				_, _, err := apiUser.DeleteMessage(ev.Channel, ev.TimeStamp)
				if err != nil {
					log.Printf("Deleting message failed: %v", err)
				}
			}
		}
	}
}

func messageMux(strArr []string, channel string) {
	reply := ""
	emoji := ":female-office-worker:"
	switch strArr[0] {
	case "cons":
		if len(strArr) < 2 {
			reply = "Please specify a sensor name: *cons <sensor>*"
		} else {
			today := time.Now().Format("06.01.02")
			cons := consumption.OneCons(strArr[1], today)
			reply = fmt.Sprintf("*%s* today's consumption: *%.2f Wh*", cons.Device, cons.Watt)
		}
	case "covid":
		sendCovid(channel)
	case "hum":
		sendHumidity(channel)
	case "temp":
		sendTemperature(channel)
	case "turn":
		if len(strArr) > 2 {
			re := regexp.MustCompile(`(?i)o(n|ff)`) // Bot name
			match := re.Match([]byte(strArr[2]))
			if !match {
				reply = "Please specify a sensor name: turn <sensor> *<on/off>*"
			}
			state := strings.ToLower(strArr[2])
			ok := setHassioService("switch."+strArr[1], "switch", "turn_"+state)
			if ok {
				reply = fmt.Sprintf("%s is successfully turned to %s", strArr[1], state)
			} else {
				reply = "Couldn't set state of sensor."
			}
			emoji = ":electric_plug:"
		}
	case "hautils":
		_, err := exec.Command("/bin/systemctl", "is-active", "--quiet", "hautils.service").Output()
		if err != nil {
			reply = "Ha utils is not running"
		} else {
			reply = "Ha utils is running"
		}

	case "help":
		reply = `
*covid* - current covid data
*cons <sensor>* - sensor consumption
*hum* - display humidity
*temp* - display temperature
*turn <sensor> <on/off>* - turn switch on/off`
	default:
		reply = fmt.Sprintf("Sorry, I dont understand \"_%s_\"", strings.Join(strArr, " "))
	}

	_, _, err := apiBot.PostMessage(channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
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
	resp := queryHassio("http://192.168.1.27:8123/api/states/"+sensor, "GET", nil)

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

func setHassioService(sensor, domain, service string) bool {
	payload := fmt.Sprintf(`{"entity_id":"%s"}`, sensor)
	link := "http://192.168.1.27:8123/api/services/" + domain + "/" + service

	resp := queryHassio(link, "POST", strings.NewReader(payload))

	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func queryHassio(url, method string, payload io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("HASS_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println(err)
	}

	return resp
}

func sendBlockMessages(channel string) {
	tempBtnText := slack.NewTextBlockObject("plain_text", "Temperature", false, false)
	humBtnText := slack.NewTextBlockObject("plain_text", "Humidity", false, false)
	covidBtnText := slack.NewTextBlockObject("plain_text", "Covid", false, false)
	tempBtn := slack.NewButtonBlockElement("", "temp", tempBtnText)
	humBtn := slack.NewButtonBlockElement("", "hum", humBtnText)
	covidBtn := slack.NewButtonBlockElement("", "covid", covidBtnText)
	actionBlock := slack.NewActionBlock("", tempBtn, humBtn, covidBtn)

	_, _, err := apiBot.PostMessage(channel, slack.MsgOptionBlocks(actionBlock))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}

func sendHumidity(channel string) {
	var sb strings.Builder
	hums := GetAllHumidity()
	for _, sensor := range hums {
		sb.WriteString(fmt.Sprintf("*%s*: %s %%\n", sensor.Name, sensor.Value))
	}
	emoji := ":droplet:"
	reply := sb.String()
	_, _, err := apiBot.PostMessage(channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}

func sendTemperature(channel string) {
	var sb strings.Builder
	hums := GetAllTemp()
	for _, sensor := range hums {
		sb.WriteString(fmt.Sprintf("*%s*: %s °C\n", sensor.Name, sensor.Value))
	}
	emoji := ":thermometer:"
	reply := sb.String()
	_, _, err := apiBot.PostMessage(channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}

func sendCovid(channel string) {
	infected, dead, cured := getCovidData()
	reply := fmt.Sprintf("*COVID*\n:biohazard_sign: *%d*\n:skull: *%d*\n:heartpulse: *%d*", infected, dead, cured)
	emoji := ":mask:"
	_, _, err := apiBot.PostMessage(channel, slack.MsgOptionText(reply, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}
