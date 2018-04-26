package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var config Config

type Config struct {
	DebugEnabled  bool
	Slack         *Slack
	SecurityToken string
	TLS           *TLS
	isInitialized bool
}

type TLS struct {
	CertFile string
	KeyFile  string
}

type Slack struct {
	Channel   string
	UserName  string
	EmojiIcon string
	Webhook   string
}

func Cfg() Config {
	if config.isInitialized == false {
		log.Panic("config file has not been loaded!")
	}

	return config
}

func LoadConfig(path string) {
	log.Println("Begin Loading Config File")

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing: ", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}

	config.isInitialized = true

	log.Println("Finished Loading Config File")
}

func Log(obj interface{}) {
	if config.DebugEnabled {
		log.Println(obj)
	}
}
