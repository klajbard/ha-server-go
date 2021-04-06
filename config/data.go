package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/klajbard/ha-server-go/types"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Arukereso     []AKItem     `yaml:"arukereso"`
	Hassio        HassioConfig `yaml:"hassio"`
	ScraperConfig string       `yaml:"scraperconfig"`
}

type AKItem struct {
	Name string   `yaml:"name"`
	Urls []string `yaml:"urls"`
}

type HassioConfig struct {
	Sensors []SensorList `yaml:"sensors"`
}

type SensorList struct {
	Type string         `yaml:"type"`
	List []SensorConfig `yaml:"list"`
}

type SensorConfig struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

var cfg = flag.String("cfg", "config.yaml", "config file path")

func (c *Configuration) GetConf() *Configuration {
	var conf []byte
	resp, err := http.Get(fmt.Sprintf("https://%s@cdn.klajbar.com/conf/ha-server-config.yaml", os.Getenv("CDN_CRED")))
	if err != nil {
		log.Println(err)
		flag.Parse()
		conf, err = ioutil.ReadFile(*cfg)
		if err != nil {
			log.Println(err)
		}
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		configData := types.AWSResponse{}

		err = json.Unmarshal([]byte(string(respBody)), &configData)
		if err != nil {
			log.Fatal(err)
		}
		conf = configData.Body.Data
	}

	if err := yaml.Unmarshal(conf, c); err != nil {
		log.Println(err)
	}

	return c
}

var Conf Configuration

func init() {
	Conf.GetConf()
}
