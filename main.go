package main

import (
	nadeo "github.com/codecat/gonadeo"
	"github.com/sirupsen/logrus"
)

var gServices nadeo.Nadeo
var gKeepAlive bool = true

func main() {
	loadConfig()

	gServices = nadeo.NewNadeoWithAudience("NadeoServices")
	gServices.SetUserAgent("Openplanet Bot / @miss / miss@openplanet.dev")

	var err error
	if appConfig.NadeoServices.Username != "" {
		err = gServices.AuthenticateBasic(appConfig.NadeoServices.Username, appConfig.NadeoServices.Password)
	} else {
		err = gServices.AuthenticateUbi(appConfig.NadeoServices.Email, appConfig.NadeoServices.Password)
	}
	if err != nil {
		logrus.WithError(err).Error("Unable to authenticate with Nadeo services")
	}

	if err := discordOpen(); err != nil {
		logrus.WithError(err).Error("Unable to initialize Discord")
		return
	}

	for _, title := range appConfig.NadeoTitles {
		go checkerTitle(title)
	}

	for _, remote := range appConfig.Remotes {
		go checkerRemote(remote)
	}

	discordClose()

	saveConfig()
}
