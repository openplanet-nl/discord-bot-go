package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/codecat/go-libs/log"
)

func checkReportRemote(info *configRemoteInfo) {
	log.Info("⚠ New update! %s = %s", info.Name, info.LastModified)

	for _, channelID := range info.Channels {
		line := fmt.Sprintf(
			":warning: **New update!** %s is now at *%s*: <%s>",
			info.Name,
			info.LastModified,
			info.URL,
		)

		_, err := appDiscord.ChannelMessageSend(channelID, line)
		if err != nil {
			log.Warn("Unable to send message to channel with ID %s", channelID)
		}
	}
}

func checkRemote(info *configRemoteInfo) bool {
	lastKnownModified := time.Time{}
	if info.LastModified != "" {
		var err error
		lastKnownModified, err = time.Parse(time.RFC3339, info.LastModified)
		if err != nil {
			log.Warn("Invalid date format: %s", err.Error())
		}
	}

	req, err := http.NewRequest("HEAD", info.URL, nil)
	if err != nil {
		log.Error("Error creating http request to %s: %s", info.Name, err.Error())
		return false
	}

	req.Header.Set("User-Agent", "Openplanet Bot v3 / @Miss#8888 / Openplanet.nl")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Warn("Error sending http request to %s: %s", info.Name, err.Error())
		return false
	}

	lastModified, err := time.Parse(time.RFC1123, res.Header.Get("Last-Modified"))
	if err != nil {
		log.Error("Invalid date format from server: %s", err.Error())
		return false
	}

	if lastModified.After(lastKnownModified) {
		info.LastModified = lastModified.Format(time.RFC3339)
		checkReportRemote(info)
		return true
	}

	return false
}

func checkerRemote(info *configRemoteInfo) {
	for gKeepAlive {
		if checkRemote(info) {
			saveConfig()
		}
		time.Sleep(time.Duration(59+rand.Intn(2)) * time.Second)
	}
}
