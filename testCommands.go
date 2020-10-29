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
	_, err := discord.ChannelMessageSend(channel, msg)
	if err != nil {
		fmt.Println(err)
	}
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
	_, err = discord.ChannelMessageSend(channel, msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg)
}

// lists all subreddits in cache along with the number of posts
func quickMemeGetCache(discord *discordgo.Session, channel string) {
	var postCount int
	var cachedReddits []string
	var cachedRedditCount int
	_, err := discord.ChannelMessageSend(channel, "Getting cache info...")
	if err != nil {
		fmt.Println(err)
	}
	cachedRedditCount = len(PostCache)
	for k := range PostCache {
		cachedReddits = append(cachedReddits, k)
		postCount += len(PostCache[k])
	}
	msgone := "There are " + strconv.Itoa(cachedRedditCount) + " cached subreddits and " + strconv.Itoa(postCount) + " posts cached."
	msgtwo := "Cached subs: " + strings.Join(cachedReddits, ", ")
	fmt.Println(msgone)
	fmt.Println(msgtwo)
	_, err = discord.ChannelMessageSend(channel, msgone)
	if err != nil {
		fmt.Println(err)
	}

	_, err = discord.ChannelMessageSend(channel, msgtwo)
	if err != nil {
		fmt.Println(err)
	}
}

// clears the cache and repopulates
func quickMemeClearCache(discord *discordgo.Session, channel string) {
	var err error
	if !CachePopulating {
		CachePopulating = true
		_, err = discord.ChannelMessageSend(channel, "Clearing cache...")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Admin issued cache clear...")
		ClearCache()
		fmt.Println("Cache cleared. Repopulating...")
		_, err = discord.ChannelMessageSend(channel, "Cache cleared. Repopulating...")
		if err != nil {
			fmt.Println(err)
		}
		st := GetMillis()
		PopulateCache()
		et := GetMillis()
		msg := "New cache time is " + strconv.FormatInt(CacheTime, 10)
		msgtwo := ". Took " + strconv.FormatInt(et-st, 10) + "ms to populate cache."
		_, err = discord.ChannelMessageSend(channel, "Done. "+msg+msgtwo)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = discord.ChannelMessageSend(channel, "Cache is currently repopulating! Please try again in a few minutes.")
		if err != nil {
			fmt.Println(err)
		}
	}
}

// gets the number of channel names cached
func quickMemeServerCache(discord *discordgo.Session, channel string) {
	//channelCount := len(ServerMap)
	channelCount := 0
	fmt.Println("There are " + strconv.Itoa(channelCount) + " text channels currently cached.")
	_, err := discord.ChannelMessageSend(channel, "There are "+strconv.Itoa(channelCount)+" text channels currently cached.")
	if err != nil {
		fmt.Println(err)
	}
}

// tests redis response time
func quickMemeTestRedis(discord *discordgo.Session, channel string) {
	_, err := discord.ChannelMessageSend(channel, "Testing Redis response time....")
	if err != nil {
		fmt.Println(err)
	}

	var times []int64
	var totalTime int64
	var avgTime float64
	for i := 0; i < 10; i++ {
		st := GetMillis()
		_, _ = GetAllBannedSubs(channel)
		et := GetMillis()
		t := et - st
		times = append(times, t)
	}
	for i := 0; i < len(times); i++ {
		totalTime += times[i]
	}
	avgTime = float64(totalTime) / float64(len(times))
	fmt.Println("Average redis response time: " + strconv.FormatFloat(avgTime, 'f', 6, 64) + " ms")
	_, err = discord.ChannelMessageSend(channel, "Average redis response time: "+strconv.FormatFloat(avgTime, 'f', 3, 64)+" ms")
	if err != nil {
		fmt.Println(err)
	}
}

// default case for the quickmeme command
func quickMemeDefault(discord *discordgo.Session, channel string) {
	servers := discord.State.Guilds
	//userCount := getNumberOfUsers(discord)
	msg := "Discord-Quick-Meme is active and ready on " + strconv.Itoa(len(servers)) + " servers."
	fmt.Println(msg)
	_, err := discord.ChannelMessageSend(channel, msg)
	if err != nil {
		fmt.Println(err)
	}
}
