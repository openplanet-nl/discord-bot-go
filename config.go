package main

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

type configRemoteInfo struct {
	Name     string
	Channels []string
	URL      string

	LastModified string
}

type configNadeoTitle struct {
	Name     string
	Channels []string
	ID       string

	Timestamp string
}

type configDiscord struct {
	Token string
}

type configNadeoServices struct {
	Email    string
	Username string
	Password string
}

type configData struct {
	Discord       configDiscord
	NadeoServices configNadeoServices

	NadeoTitles []*configNadeoTitle
	Remotes     []*configRemoteInfo
}

var appConfig configData

func loadConfig() {
	configBytes, err := os.ReadFile("config.toml")
	if err != nil {
		logrus.WithError(err).Error("Unable to read config.toml file")
		return
	}

	err = toml.Unmarshal(configBytes, &appConfig)
	if err != nil {
		logrus.WithError(err).Error("Unable to unmarshal config.toml")
		return
	}
}

func saveConfig() {
	f, err := os.Create("config.toml")
	if err != nil {
		logrus.WithError(err).Error("Unable to open config file for writing")
		return
	}

	configBytes, err := toml.Marshal(&appConfig)
	if err != nil {
		logrus.WithError(err).Error("Unable to marshal config to file")
		return
	}

	err = os.WriteFile("config.toml", configBytes, 0644)
	if err != nil {
		logrus.WithError(err).Error("Unable to write to config file")
		return
	}

	f.Close()
}
