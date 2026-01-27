package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func checkReportRemote(info *configRemoteInfo) {
	logrus.WithFields(logrus.Fields{
		"name":          info.Name,
		"last-modified": info.LastModified,
	}).Info("⚠ New update!")

	for _, channelID := range info.Channels {
		line := fmt.Sprintf(
			":warning: **New update!** %s is now at *%s*: <%s>",
			info.Name,
			info.LastModified,
			info.URL,
		)

		_, err := appDiscord.ChannelMessageSend(channelID, line)
		if err != nil {
			logrus.Warn("Unable to send message to channel with ID ", channelID)
		}
	}
}

func checkRemote(info *configRemoteInfo) bool {
	lastKnownModified := time.Time{}
	if info.LastModified != "" {
		var err error
		lastKnownModified, err = time.Parse(time.RFC3339, info.LastModified)
		if err != nil {
			logrus.WithError(err).Warn("Invalid date format")
		}
	}

	req, err := http.NewRequest("HEAD", info.URL, nil)
	if err != nil {
		logrus.WithError(err).Error("Error creating http request to ", info.Name)
		return false
	}

	req.Header.Set("User-Agent", "Openplanet Bot / @miss / miss@openplanet.dev")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Warn("Error sending http request to ", info.Name)
		return false
	}

	lastModified, err := time.Parse(time.RFC1123, res.Header.Get("Last-Modified"))
	if err != nil {
		logrus.WithError(err).Error("Invalid date format from server")
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
