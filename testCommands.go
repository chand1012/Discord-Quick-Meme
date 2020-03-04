package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// does a quick speed test of reddit
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

// lists all subreddits in cache along with the number of posts
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

// clears the cache and repopulates
func quickMemeClearCache(discord *discordgo.Session, channel string) {
	if !CachePopulating {
		CachePopulating = true
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
	} else {
		discord.ChannelMessageSend(channel, "Cache is currently repopulating! Please try again in a few minutes.")
	}
}

// gets the number of channel names cached
func quickMemeServerCache(discord *discordgo.Session, channel string) {
	//channelCount := len(ServerMap)
	channelCount := 0
	fmt.Println("There are " + strconv.Itoa(channelCount) + " text channels currently cached.")
	discord.ChannelMessageSend(channel, "There are "+strconv.Itoa(channelCount)+" text channels currently cached.")
}

// tests redis response time
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

// gets the number of common subreddits to cache
func quickMemeTestCommonCache(discord *discordgo.Session, channel string) {
	discord.ChannelMessageSend(channel, "Getting stats on Common Subreddits....")
	discord.ChannelMessageSend(channel, "There are "+strconv.Itoa(int(CommonSubsCounter))+" additional subs being cached.")
	values, _ := getCommonSubsRedisRaw()
	discord.ChannelMessageSend(channel, "Subs stored in Redis cache: "+values)
	discord.ChannelMessageSend(channel, "Getting times for the cached subreddits....")
	var sendStr string
	counter := 0
	sendStr = ""
	for key := range CommonSubs {
		sendStr += (strconv.Itoa(counter) + ") " + key + ": Count: " + strconv.Itoa(int(CommonSubs[key])) + " Time: " + strconv.FormatInt(CommonSubsTime[key], 10) + " ms\n")
		counter++
	}
	discord.ChannelMessageSend(channel, sendStr)
}

// default case for the quickmeme command
func quickMemeDefault(discord *discordgo.Session, channel string) {
	servers := discord.State.Guilds
	userCount := getNumberOfUsers(discord)
	msg := "Discord-Quick-Meme is active and ready on " + strconv.Itoa(len(servers)) + " servers for " + strconv.Itoa(userCount) + " users."
	fmt.Println(msg)
	discord.ChannelMessageSend(channel, msg)
}
