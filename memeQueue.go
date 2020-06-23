package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// QueueObj data structure for the queue
type QueueObj struct {
	Interval   string   `json:"interval"`
	Time       int64    `json:"time"`
	Type       string   `json:"type"`
	SubReddits []string `json:"subreddit"`
	NSFW       bool     `json:"nsfw"`
}

func queueThread(discord *discordgo.Session) {

	fmt.Println("Starting Queue Processing thread.")
	fmt.Println("Generating lock file.")
	testData, err := lockFileCreate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Thread started.")
	for {

		var wg sync.WaitGroup

		fileEqual, err := lockFileEqu(testData)

		if err != nil {
			fmt.Println(err)
			break
		}

		if !fileEqual {
			fmt.Println("New processing thread started, killing old thread.")
			break
		}

		keys, err := GetAllQueueChannels()
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
		// This uses less CPU than the loop with timer
		// And the end user will see no difference
		// One of the few times I will use a sleep
		time.Sleep(time.Second * 10)
	}
	fmt.Println("Queue processing thread killed.")
}

func queueWorker(discord *discordgo.Session, channel string, wg *sync.WaitGroup) {

	defer wg.Done()

	var interval time.Duration

	maxTime := time.Hour * 168
	minTime := time.Minute * 15

	queueItem, err := GetMemeQueue(channel)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		fmt.Println("Error getting from Queue: ", err.Error())
	}

	if queueItem.Time <= time.Now().Unix() && !QueueState[channel] {

		if err != nil {
			// will only send the error to the channel if its time
			// for the channel to get a meme. This should happen infrequently
			// or never
			errSendRoutine(discord, channel, err)
			return
		}

		QueueState[channel] = true

		// will always set the QueueState for the
		// given channel to false no matter where the
		// function ends.
		defer resetQueueState(channel)

		fmt.Println("Posting in", channel, "from queue.")

		switch queueItem.Type {
		case "text":
			getTextPost(discord, channel, queueItem.NSFW, queueItem.SubReddits, "hot")
		case "link":
			getLinkPost(discord, channel, queueItem.NSFW, queueItem.SubReddits, "hot")
		default:
			getMediaPost(discord, channel, queueItem.NSFW, queueItem.SubReddits, "hot")
		}
		letters, err := regexp.Compile("[^a-zA-Z]+")

		if err != nil {
			fmt.Println("Error setting Queue: ", err.Error())
			errSendRoutine(discord, channel, err)
			return
		}

		numbers, err := regexp.Compile("[^0-9]+")

		if err != nil {
			fmt.Println("Error setting Queue: ", err.Error())
			errSendRoutine(discord, channel, err)
			return
		}

		multiplier, err := strconv.ParseFloat(numbers.ReplaceAllString(queueItem.Interval, ""), 64)

		if err != nil {
			fmt.Println("Error setting Queue: ", err.Error())
			errSendRoutine(discord, channel, err)
			return
		}

		duration := letters.ReplaceAllString(queueItem.Interval, "")

		if err != nil {
			fmt.Println("Error setting Queue: ", err.Error())
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

		if interval > maxTime {
			interval = maxTime
			queueItem.Interval = "1w"
			discord.ChannelMessageSend(channel, "The maximum interval between memes is one week, so the interval will be set to that. For slower memes, check Facebook or a newspaper.")
		}

		if interval < minTime {
			interval = minTime
			queueItem.Interval = "15m"
			discord.ChannelMessageSend(channel, "There is minimum time interval between posts of 15 minutes, setting the interval to that. These servers aren't free, you know.")
		}

		queueItem.Time = time.Now().Unix() + int64(interval.Seconds())

		err = UpdateMemeQueueTime(channel, queueItem.Time)

		if err != nil {
			fmt.Println("Error setting Queue: ", err.Error())
			errSendRoutine(discord, channel, err)
		} else {
			fmt.Println("Done.")
		}
	}
}

func lockFileEqu(input []byte) (bool, error) {
	data, err := ioutil.ReadFile("./thread.lock")
	if err != nil {
		return false, err
	}
	if bytes.Compare(input, data) == 0 {
		return true, nil
	}
	return false, nil
}

func lockFileExists() bool {
	info, err := os.Stat("./thread.lock")
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func lockFileCreate() ([]byte, error) {
	fileData := make([]byte, 8)
	rand.Read(fileData)
	err := ioutil.WriteFile("./thread.lock", fileData, 0644)
	return fileData, err
}

func resetQueueState(channel string) {
	QueueState[channel] = false
}
