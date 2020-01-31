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

func quickMemeImageSearch(discord *discordgo.Session, channel string, imageURL string) {
	discord.ChannelMessageSend(channel, "Searching the web...")
	lastMessages := discord.ChannelMessages(channel, 100, "", "", "")

}
