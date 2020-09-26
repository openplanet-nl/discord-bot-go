package main

import (
	"github.com/codecat/go-libs/log"
)

func main() {
	loadConfig()

	if err := discordOpen(); err != nil {
		log.Error("Unable to initialize Discord: %s", err.Error())
		return
	}

	for _, remote := range appConfig.Remotes {
		go checker(remote)
	}

	discordClose()

	saveConfig()
}
