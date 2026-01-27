package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var appDiscord *discordgo.Session

func discordReady(s *discordgo.Session, event *discordgo.Ready) {
	logrus.Info("Discord connected: ", appDiscord.State.User)
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	gKeepAlive = false
	appDiscord.Close()
}
