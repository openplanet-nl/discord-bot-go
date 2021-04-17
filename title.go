package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/codecat/go-libs/log"
)

type encryptedPackage struct {
	EncryptedPackageID              string `json:"encryptedPackageId"`
	Name                            string
	PublicEncryptedPackageVersionID string `json:"publicEncryptedPackageVersionId"`
	Timestamp                       string
	TitleID                         string `json:"titleId"`
}

func checkReportTitle(title *configNadeoTitle) {
	log.Info("⚠ New update! %s = %s", title.Name, title.Timestamp)

	for _, channelID := range title.Channels {
		line := fmt.Sprintf(
			":warning: **New update!** %s is now at *%s*",
			title.Name,
			title.Timestamp,
		)

		_, err := appDiscord.ChannelMessageSend(channelID, line)
		if err != nil {
			log.Warn("Unable to send message to channel with ID %s", channelID)
		}
	}
}

func checkTitle(title *configNadeoTitle) bool {
	res, err := gServices.Get(
		fmt.Sprintf(
			"https://prod.trackmania.core.nadeo.online/encryptedPackages/?titleIdList=%s",
			title.ID,
		),
		false,
	)

	if err != nil {
		log.Error("Unable to check title updates: %s", err.Error())
		return false
	}

	var packages []encryptedPackage
	json.Unmarshal([]byte(res), &packages)

	if len(packages) == 0 {
		log.Error("Title does not exist: %s", title.ID)
		return false
	}

	tm, err := time.Parse(time.RFC3339, packages[0].Timestamp)
	if err != nil {
		log.Error("Unable to parse timestamp: %s", err.Error())
		return false
	}

	lastKnownTimestamp := time.Time{}
	if title.Timestamp != "" {
		lastKnownTimestamp, _ = time.Parse(time.RFC3339, title.Timestamp)
	}

	if tm.After(lastKnownTimestamp) {
		title.Timestamp = tm.Format(time.RFC3339)
		checkReportTitle(title)
		return true
	}

	return false
}

func checkerTitle(title *configNadeoTitle) {
	for gKeepAlive {
		if checkTitle(title) {
			saveConfig()
		}
		time.Sleep(time.Duration(59+rand.Intn(2)) * time.Second)
	}
}
