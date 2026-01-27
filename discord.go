package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var gDiscord *discordgo.Session

func discordReady(s *discordgo.Session, event *discordgo.Ready) {
	logrus.Info("Discord connected: ", gDiscord.State.User)
}

func discordOpen() error {
	var err error
	gDiscord, err = discordgo.New("Bot " + appConfig.Discord.Token)
	if err != nil {
		return err
	}

	gDiscord.AddHandler(discordReady)

	err = gDiscord.Open()
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
	gDiscord.Close()
}
