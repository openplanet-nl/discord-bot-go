package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/codecat/go-libs/log"
)

var appDiscord *discordgo.Session

func discordReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Discord connected: %s", appDiscord.State.User)
}

func discordOpen() error {
	var err error
	appDiscord, err = discordgo.New("Bot " + appConfig.Discord.Token)
	if err != nil {
		return err
	}

	appDiscord.AddHandler(discordReady)

	err = appDiscord.Open()
	if err != nil {
		return err
	}

	return nil
}

func discordClose() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	gKeepAlive = false
	appDiscord.Close()
}
