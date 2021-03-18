package hass

import (
	"flag"
	"io/ioutil"
	"log"

	"../config"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Marketplace  MarketplaceConfig `yaml:"marketplace"`
	HaBump       []HaBumpConfig    `yaml:"habump"`
	Channels     []SlackChannel    `yaml:"channels"`
	StockWatcher []ItemStock       `yaml:"stockwatcher"`
	Arukereso    []Url             `yaml:"arukereso"`
	Dht          DhtConfig         `yaml:"dht"`
	Enable       EnableConfig      `yaml:"enable"`
	Silence      bool              `yaml:"silence"`
}

type Url struct {
	Url string `yaml:"url"`
}

type MarketplaceConfig struct {
	Jofogas     []MarketplaceName `yaml:"jofogas"`
	Hardverapro []MarketplaceName `yaml:"hardverapro"`
}

type HaBumpConfig struct {
	Identifier string   `yaml:"identifier"`
	Items      []HaItem `yaml:"items"`
}

type DhtConfig struct {
	Pin int `yaml:"pin"`
}

type HaItem struct {
	Name  string `yaml:"name"`
	Id    string `yaml:"id"`
	Hour  int    `yaml:"hour"`
	Start int    `yaml:"start"`
}

type MarketplaceName struct {
	Name string `yaml:"name"`
}

type SlackChannel struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type ItemStock struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Query string `yaml:"query"`
}

type EnableConfig struct {
	Bestbuy      bool `yaml:"bestbuy"`
	Stockwatcher bool `yaml:"stockwatcher"`
	Marketplace  bool `yaml:"marketplace"`
	Steamgifts   bool `yaml:"steamgifts"`
	Dht          bool `yaml:"dht"`
	Arukereso    bool `yaml:"arukereso"`
	Covid        bool `yaml:"covid"`
	Bumphva      bool `yaml:"bumphva"`
	Ncore        bool `yaml:"ncore"`
	Fuel         bool `yaml:"fuel"`
	Fixerio      bool `yaml:"fixerio"`
	Awscost      bool `yaml:"awscost"`
	Btc          bool `yaml:"btc"`
}

var cfg = flag.String("hagoconf", config.Conf.ScraperConfig, "config file path")

func Get() *Configuration {
	var c *Configuration
	flag.Parse()

	conf, err := ioutil.ReadFile(*cfg)
	if err != nil {
		log.Println(err)
	}
	if err := yaml.Unmarshal(conf, &c); err != nil {
		log.Println(err)
	}

	return c
}
