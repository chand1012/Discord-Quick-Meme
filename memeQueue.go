package main

import (
	"fmt"
	"regexp"
	"strconv"
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

	defer wg.Done()

	var interval time.Duration

	maxTime := time.Hour * 168

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

		letters, err := regexp.Compile("[^a-zA-Z]+")

		if err != nil {
			fmt.Println("Error setting redis queue: ", err.Error())
			errSendRoutine(discord, channel, err)
			return
		}

		numbers, err := regexp.Compile("[^0-9]+")

		if err != nil {
			fmt.Println("Error setting redis queue: ", err.Error())
			errSendRoutine(discord, channel, err)
			return
		}

		multiplier, err := strconv.ParseFloat(numbers.ReplaceAllString(queueItem.Interval, ""), 64)

		if err != nil {
			fmt.Println("Error setting redis queue: ", err.Error())
			errSendRoutine(discord, channel, err)
			return
		}

		duration := letters.ReplaceAllString(queueItem.Interval, "")

		if err != nil {
			fmt.Println("Error setting redis queue: ", err.Error())
			errSendRoutine(discord, channel, err)
			return
		}

		if multiplier <= 0 {
			multiplier = 1
			discord.ChannelMessageSend(channel, "Time cannot be negative or zero, setting to one.")
		}

		switch duration {
		case "s":
			interval = time.Second * time.Duration(multiplier)
		case "m":
			interval = time.Minute * time.Duration(multiplier)
		case "d":
			interval = time.Hour * 24 * time.Duration(multiplier)
		case "w":
			interval = time.Hour * 168 * time.Duration(multiplier)
		default:
			interval = time.Hour * time.Duration(multiplier)
		}

		queueItem.Time = interval.Milliseconds()

		if queueItem.Time > maxTime.Milliseconds() {
			queueItem.Time = maxTime.Milliseconds()
			discord.ChannelMessageSend(channel, "The maximum interval between memes is one week, so the interval will be set to that. For slower memes, check Facebook or a newspaper.")
		}

		err = setRedisQueueRaw(queueItem, channel)

		if err != nil {
			fmt.Println("Error setting redis queue: ", err.Error())
			errSendRoutine(discord, channel, err)
		}
	}
}
