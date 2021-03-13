package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"../consumption"
	"../types"
	"../utils"
)

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

func TurnSwitch(strArr []string, channel string) {
	reply := "Wrong parameters"
	emoji := ":electric_plug:"
	if len(strArr) < 2 {
		utils.PostMessage(channel, reply, emoji)
	}
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
	utils.PostMessage(channel, reply, emoji)
}

func Humidity(channel string) {
	var sb strings.Builder
	hums := getAllHumidity()
	for _, sensor := range hums {
		sb.WriteString(fmt.Sprintf("*%s*: %s %%\n", sensor.Name, sensor.Value))
	}
	emoji := ":droplet:"
	reply := sb.String()

	utils.PostMessage(channel, reply, emoji)
}

func Temperature(channel string) {
	var sb strings.Builder
	hums := getAllTemp()
	for _, sensor := range hums {
		sb.WriteString(fmt.Sprintf("*%s*: %s °C\n", sensor.Name, sensor.Value))
	}
	emoji := ":thermometer:"
	reply := sb.String()

	utils.PostMessage(channel, reply, emoji)
}

func Consumption(strArr []string, channel string) {
	emoji := ":electric_plug:"
	reply := "Please specify a sensor name: *cons <sensor>*"
	if len(strArr) > 1 {
		today := time.Now().Format("06.01.02")
		cons := consumption.OneCons(strArr[1], today)
		reply = fmt.Sprintf("*%s* today's consumption: *%.2f Wh*", cons.Device, cons.Watt)
	}

	utils.PostMessage(channel, reply, emoji)
}

func getAllTemp() []types.SensorValue {
	ret := []types.SensorValue{}
	for sensor, name := range TEMPERATURES {
		ret = append(ret, types.SensorValue{Name: name, Value: getHassioData(sensor)})
	}

	return ret
}

func getAllHumidity() []types.SensorValue {
	ret := []types.SensorValue{}
	for sensor, name := range HUMIDITIES {
		ret = append(ret, types.SensorValue{Name: name, Value: getHassioData(sensor)})
	}

	return ret
}

func getHassioData(sensor string) string {
	sensorData := &types.SensorData{}
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
