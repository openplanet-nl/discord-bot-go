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
	if err := gServices.AuthenticateUbi(appConfig.NadeoServices.Email, appConfig.NadeoServices.Password); err != nil {
		logrus.WithError(err).Error("Unable to authenticate with Nadeo services")
		return
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
