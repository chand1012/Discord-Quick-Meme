package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func queueWorker(discord *discordgo.Session) {
	fmt.Println("Starting Queue Processing thread.")
	for {
		keys, err := getAllQueueChannels()

		for _, key := range keys {
			queueItem, err := getRedisQueue(key)

		}
	}
}
