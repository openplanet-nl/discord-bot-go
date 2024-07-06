package main

import (
	"github.com/codecat/go-libs/log"
	nadeo "github.com/codecat/gonadeo"
)

var gServices nadeo.Nadeo
var gKeepAlive bool = true

func main() {
	loadConfig()

	gServices = nadeo.NewNadeoWithAudience("NadeoServices")
	gServices.SetUserAgent("Openplanet Bot / @miss / miss@openplanet.dev")
	if err := gServices.AuthenticateUbi(appConfig.NadeoServices.Email, appConfig.NadeoServices.Password); err != nil {
		log.Error("Unable to authenticate with Nadeo services: %s", err.Error())
		return
	}

	if err := discordOpen(); err != nil {
		log.Error("Unable to initialize Discord: %s", err.Error())
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
