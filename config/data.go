package config

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Arukereso []AKItem     `yaml:"arukereso"`
	Hassio    HassioConfig `yaml:"hassio"`
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

func (c *Configuration) GetConf() *Configuration {
	cfg := flag.String("cfg", "config.yaml", "config file path")
	flag.Parse()

	conf, err := ioutil.ReadFile(*cfg)
	if err != nil {
		log.Println(err)
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
