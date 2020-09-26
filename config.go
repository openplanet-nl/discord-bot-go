package main

import (
	"io/ioutil"
	"os"

	"github.com/codecat/go-libs/log"
	"github.com/pelletier/go-toml"
)

type configRemoteInfo struct {
	Name     string
	URL      string
	Channels []string

	LastModified string
}

type configDiscord struct {
	Token string
}

type configData struct {
	Discord configDiscord
	Remotes []*configRemoteInfo
}

var appConfig configData

func loadConfig() {
	configBytes, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Error("Unable to read config.toml file: %s", err.Error())
		return
	}

	err = toml.Unmarshal(configBytes, &appConfig)
	if err != nil {
		log.Error("Unable to unmarshal config.toml: %s", err.Error())
		return
	}
}

func saveConfig() {
	f, err := os.Create("config.toml")
	if err != nil {
		log.Error("Unable to open config file for writing: %s", err.Error())
		return
	}

	configBytes, err := toml.Marshal(&appConfig)
	if err != nil {
		log.Error("Unable to marshal config to file: %s", err.Error())
		return
	}

	err = ioutil.WriteFile("config.toml", configBytes, 0644)
	if err != nil {
		log.Error("Unable to write to config file: %s", err.Error())
		return
	}

	f.Close()
}
