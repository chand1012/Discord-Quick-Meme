package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// this is for all of the commands that have a lot of extra crap in them

func quickMemeTest(discord *discordgo.Session, channel string) {
	var count int
	var total int64
	var redditResult float64
	var msg string
	msg = "Starting Quick-Meme speed test..."
	discord.ChannelMessageSend(channel, msg)
	fmt.Println(msg)
	for i := 0; i < 10; i++ {
		starttime := GetMillis()
		_ = PingReddit()
		endtime := GetMillis()
		total += (endtime - starttime)
		count = i
	}
	redditResult = float64(total) / float64(count)
	msg = "Average Reddit response time over 10 trials: " + strconv.FormatFloat(redditResult, 'f', 1, 64) + "ms"
	discord.ChannelMessageSend(channel, msg)
	fmt.Println(msg)
}

func quickMemeGetCache(discord *discordgo.Session, channel string) {
	var postCount int
	var cachedReddits []string
	var cachedRedditCount int
	discord.ChannelMessageSend(channel, "Getting cache info...")
	cachedRedditCount = len(PostCache)
	for k := range PostCache {
		cachedReddits = append(cachedReddits, k)
		postCount += len(PostCache[k])
	}
	msgone := "There are " + strconv.Itoa(cachedRedditCount) + " cached subreddits and " + strconv.Itoa(postCount) + " posts cached."
	msgtwo := "Cached subs: " + strings.Join(cachedReddits, ", ")
	fmt.Println(msgone)
	fmt.Println(msgtwo)
	discord.ChannelMessageSend(channel, msgone)
	discord.ChannelMessageSend(channel, msgtwo)
}

func quickMemeClearCache(discord *discordgo.Session, channel string) {
	discord.ChannelMessageSend(channel, "Clearing cache...")
	fmt.Println("Admin issued cache clear...")
	ClearCache()
	fmt.Println("Cache cleared. Repopulating...")
	discord.ChannelMessageSend(channel, "Cache cleared. Repopulating...")
	st := GetMillis()
	PopulateCache()
	et := GetMillis()
	msg := "New cache time is " + strconv.FormatInt(CacheTime, 10)
	msgtwo := ". Took " + strconv.FormatInt(et-st, 10) + "ms to populate cache."
	discord.ChannelMessageSend(channel, "Done. "+msg+msgtwo)
}

func quickMemeServerCache(discord *discordgo.Session, channel string) {
	channelCount := len(ServerMap)
	fmt.Println("There are " + strconv.Itoa(channelCount) + " text channels currently cached.")
	fmt.Println(ServerMap)
	discord.ChannelMessageSend(channel, "There are "+strconv.Itoa(channelCount)+" text channels currently cached.")
}

func quickMemeTestRedis(discord *discordgo.Session, channel string) {
	discord.ChannelMessageSend(channel, "Testing Redis response time....")
	var times []int64
	var totalTime int64
	var avgTime float64
	for i := 0; i < 10; i++ {
		st := GetMillis()
		_, _ = GetBannedSubreddits(channel)
		et := GetMillis()
		t := et - st
		times = append(times, t)
	}
	for i := 0; i < len(times); i++ {
		totalTime += times[i]
	}
	avgTime = float64(totalTime) / float64(len(times))
	fmt.Println("Average redis response time: " + strconv.FormatFloat(avgTime, 'f', 6, 64) + " ms")
	discord.ChannelMessageSend(channel, "Average redis response time: "+strconv.FormatFloat(avgTime, 'f', 3, 64)+" ms")
}

func quickMemeImageSearch(discord *discordgo.Session, channel string) {
	searchURL := ""
	extensions := []string{".jpg", ".png", ".jpeg"}
	discord.ChannelMessageSend(channel, "Searching Reddit...")
	lastMessages, _ := discord.ChannelMessages(channel, 100, "", "", "")
	for _, message := range lastMessages {
		attachmentsLength := len(message.Attachments)
		currentIndex := attachmentsLength - 1
		if attachmentsLength > 0 {
			for currentIndex >= 0 {
				last := message.Attachments[currentIndex]
				if ContainsAnySubstring(last.ProxyURL, extensions) {
					searchURL = last.ProxyURL
					break
				} else {
					currentIndex--
				}
			}
		}

		if searchURL != "" {
			break
		}

		if ContainsAnySubstring(message.Content, extensions) && !message.Author.Bot {
			httpIndex := strings.Index(message.Content, "http")
			startLink := message.Content[httpIndex:len(message.Content)]
			spaceIndex := strings.Index(startLink, " ")
			if spaceIndex != -1 {
				searchURL = startLink[0:spaceIndex]
				break
			} else {
				searchURL = startLink
			}
		}
	}

	if searchURL == "" {
		discord.ChannelMessageSend(channel, "404: URL not found. If you think that this is a mistake, post on our Github issues page along with appropriate screenshots and information. https://github.com/chand1012/Discord-Quick-Meme/issues")
		return
	}

	url := imageRedditSearch(searchURL)
	if url == "" {
		discord.ChannelMessageSend(channel, "Couldn't find anything on Reddit, searching the web....")
		urls := imageSearch(searchURL)
		if urls == nil {
			discord.ChannelMessageSend(channel, "500: Error connecting to image search service. If this persists, report at the Github issues page found here: https://github.com/chand1012/Discord-Quick-Meme/issues")
			return
		}
		printstr := "Found " + strconv.Itoa(len(urls)) + " results:\n"
		for _, link := range urls {
			printstr += link + "\n"
		}
		discord.ChannelMessageSend(channel, printstr)
	} else if url == "nil" {
		discord.ChannelMessageSend(channel, "500: Error connecting to image search service. If this persists, report at the Github issues page found here: https://github.com/chand1012/Discord-Quick-Meme/issues")
	} else {
		discord.ChannelMessageSend(channel, "I think I found the meme: \n"+url)
	}
}
