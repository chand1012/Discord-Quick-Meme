package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)

func queueWorker(discord *discordgo.Session) {
	var timer int64
	var checkInterval int64
	checkInterval = 60
	timer = 0
	fmt.Println("Starting Queue Processing thread.")
	for {
		if timer <= time.Now().Unix() {
			fmt.Println("Checking queue...")
			keys, err := getAllQueueChannels()
			if err != nil {
				fmt.Println("Error running worker queue: ", err.Error())
				continue
			}
			fmt.Println(keys)
			for _, key := range keys {
				if key == "" {
					continue
				}
				var newTime int64
				fmt.Println(key)
				queueItem, err := getRedisQueue(key)
				if err == redis.Nil {
					continue
				}
				if err != nil {
					fmt.Println("Error getting from redis queue: ", err.Error())
					errSendRoutine(discord, key, err)
					continue
				}
				subs := &queueItem.SubReddits
				interval := &queueItem.Interval
				postTime := &queueItem.Time
				postType := &queueItem.Type
				postNSFW := &queueItem.NSFW
				customInterval := &queueItem.CustomInterval
				if *postTime <= time.Now().Unix() {
					switch *postType {
					case "media":
						getMediaPost(discord, key, *postNSFW, *subs, "hot")
					case "text":
						getTextPost(discord, key, *postNSFW, *subs, "hot")
					case "link":
						getLinkPost(discord, key, *postNSFW, *subs, "hot")
					}

					switch *interval {
					case "hourly":
						newTime = time.Now().Unix() + 3600
					case "daily":
						newTime = time.Now().Unix() + 86400
					case "twice daily":
						newTime = time.Now().Unix() + 43200
					case "custom":
						newTime = time.Now().Unix() + *customInterval
					}
					queueItem.Time = newTime
					err = setRedisQueueRaw(queueItem, key)
					if err != nil {
						fmt.Println("Error setting redis queue: ", err.Error())
						errSendRoutine(discord, key, err)
					}
				}
			}
			timer = time.Now().Unix() + checkInterval
		}
	}
}
