package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)

func queueThread(discord *discordgo.Session) {
	var checkInterval int64
	var timer int64
	checkInterval = 10
	timer = 0
	fmt.Println("Starting Queue Processing thread.")

	for {
		if timer <= time.Now().Unix() {
			var wg sync.WaitGroup

			keys, err := getAllQueueChannels()
			if err != nil {
				fmt.Println("Error running worker queue: ", err.Error())
				continue
			}

			for _, key := range keys {
				if key != "" {
					wg.Add(1)
					go queueWorker(discord, key, &wg)
				}
			}

			wg.Wait()
			timer = time.Now().Unix() + checkInterval
		}
	}
}

func queueWorker(discord *discordgo.Session, channel string, wg *sync.WaitGroup) {
	var newTime int64

	defer wg.Done()

	queueItem, err := getRedisQueue(channel)
	if err == redis.Nil {
		return
	}
	if err != nil {
		fmt.Println("Error getting from redis queue: ", err.Error())
		errSendRoutine(discord, channel, err)
		return
	}

	if queueItem.Time <= time.Now().Unix() && !QueueState[channel] {

		QueueState[channel] = true

		fmt.Println("Posting in", channel, "from queue.")

		switch queueItem.Type {
		case "media":
			getMediaPost(discord, channel, queueItem.NSFW, queueItem.SubReddits, "hot")
		case "text":
			getTextPost(discord, channel, queueItem.NSFW, queueItem.SubReddits, "hot")
		case "link":
			getLinkPost(discord, channel, queueItem.NSFW, queueItem.SubReddits, "hot")
		}

		switch queueItem.Interval { // this iwll be replaced with custom times only
		case "1h":
			newTime = time.Now().Unix() + 3600
		case "1d":
			newTime = time.Now().Unix() + 86400
		case "12h":
			newTime = time.Now().Unix() + 43200
		case "1w":
			newTime = time.Now().Unix() + 604800
		case "6h":
			newTime = time.Now().Unix() + 21600
		}

		queueItem.Time = newTime
		err = setRedisQueueRaw(queueItem, channel)

		if err != nil {
			fmt.Println("Error setting redis queue: ", err.Error())
			errSendRoutine(discord, channel, err)
		}
	}
}
