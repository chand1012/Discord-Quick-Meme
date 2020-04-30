package main

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func queueWorker(discord *discordgo.Session, wg *sync.WaitGroup) {
	fmt.Println("Starting Queue Processing thread.")
	for {
		keys, err := getAllQueueChannels()

	}
}
