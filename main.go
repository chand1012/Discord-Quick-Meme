package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var (
	commandPrefix string
	botID         string
	adminIDs      []string
	topgg         string
	// CacheTime stores cache timer value
	CacheTime int64
	//BlacklistTime stores the blacklist time for all of the channels
	BlacklistTime int64
	// ServerMap this is all of the servers an the servers this gets wiped from memory as soon as the Bot gets killed
	ServerMap map[string]string
	// NSFWMap stores all nsfw values for each channel
	NSFWMap map[string]bool
	//PostCache stores all posts
	PostCache map[string][]QuickPost
	//Blacklist list of all of the post that are blacklisted from the specified channel
	Blacklist map[string][]QuickPost // will be wiped every two to three hours
	//CommonSubs stores the amount of times the subs are hit
	CommonSubs map[string]uint8 // only needs to count up to 10
	//CommonSubsTime if one week passes, clear the above cache
	// 604800000ms in a week
	CommonSubsTime map[string]int64
	// CommonSubsCounter counts the number of common subs
	CommonSubsCounter uint8
	//LastPost gets the last post from the specified channel string
	LastPost map[string]QuickPost
	//SubMap contains all of the data for the subs
	SubMap map[string][]string
	//CachePopulating if true, do not run the populate cache until finished
	CachePopulating bool
	mrisaAddress    string
	//ErrorMsg Main error message that gets send when something goes seriously wrong
	ErrorMsg string
	//JSONError JSON error message
	JSONError string
	// RequestCount counts how many requests a channel makes
	// If over 10 in a minute then halt all posting
	RequestCount map[string]uint8
	// RequestTimer stores channel timers
	// Resets every minute
	RequestTimer map[string]int64
	// RunMode the mode that the bot is running in
	RunMode string
)

// main loop
func main() {
	var err error
	var file string
	var key string
	var adminRawIDs []string
	ServerMap = make(map[string]string)
	NSFWMap = make(map[string]bool)
	PostCache = make(map[string][]QuickPost)
	Blacklist = make(map[string][]QuickPost)
	CommonSubs = make(map[string]uint8)
	CommonSubsTime = make(map[string]int64)
	CommonSubsCounter = 0
	LastPost = make(map[string]QuickPost)
	SubMap = make(map[string][]string)
	RequestCount = make(map[string]uint8)
	RequestTimer = make(map[string]int64)
	RunMode = getMode("data.json")
	ErrorMsg = "There was an error processing your request. If this persists, please submit a report here: https://github.com/chand1012/Discord-Quick-Meme/issues"
	JSONError = "Error reading JSON file"
	file = "data.json"
	key, adminRawIDs, topgg, err = loginExtract(file)
	mrisaAddress = mrisaExtract(file)
	if err != nil {
		panic(err) // can't run without a login
	}
	discord, err := discordgo.New("Bot " + key)

	if err != nil {
		panic(err)
	}

	user, err := discord.User("@me")
	if err != nil {
		panic(err)
	}
	for _, admin := range adminRawIDs {
		a, err := discord.User(admin)
		if err != nil {
			panic(err)
		}
		adminIDs = append(adminIDs, a.ID)
	}
	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(readyHandler)
	err = discord.Open()
	if err != nil {
		panic(err)
	}
	defer discord.Close()
	commandPrefix = "!"
	<-make(chan struct{})
}

// handles bot initialization
func readyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	servers := discord.State.Guilds
	getAllChannelNames(discord)
	CachePopulating = true
	PopulateCache()
	ResetBlacklist()
	serverCount := int64(len(servers))
	fmt.Println("Discord-Quick-Meme has started on " + humanize.Comma(serverCount) + " servers")
	if RunMode == "prod" { // only run if production
		go updateServerCount(serverCount, topgg)
	}
	go updateStatus(discord)
	go queueThread(discord)
}

// handes incoming commands
func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	var subs []string
	var channelName string
	go updateStatus(discord)
	go UpdateBlacklistTime()
	channel := message.ChannelID
	commands := []string{"!meme", "!joke", "!hentai", "!news", "!fiftyfifty", "!5050", "!all", "!quickmeme", "!text", "!link", "!source", "!buzzword", "!revsearch"}
	user := message.Author
	content := message.Content
	commandContent := strings.Split(content, " ")
	command := commandContent[0]
	guildID := message.GuildID
	if user.ID == botID || user.Bot || !stringInSlice(commandContent[0], commands) {
		return
	}
	canUserPost := updateChannelTimer(channel)
	if !canUserPost {
		if RequestCount[channel] == 6 {
			fmt.Println("Channel " + channel + "is sending a lot of requests, limiting their input for 60 seconds.")
		}
		if RequestCount[channel] > 10 {
			discord.ChannelMessageSend(channel, "You're sending a lot of requests, how about you slow it down a bit? All requests from this channel will be ignored for 60 seconds.")
		}
		return
	}
	channelName = "#" + getChannelName(discord, channel, guildID)
	fmt.Println("Command '" + content + "' from " + user.Username + " on " + channelName + " (" + channel + ")")
	nsfw := getChannelNSFW(discord, channel, guildID)
	switch {
	case command == "!meme" && len(commandContent) == 1:
		subs = SubMap["memes"]
		getMediaPost(discord, channel, nsfw, subs, "hot")
	case command == "!meme" && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		if subs == nil {
			discord.ChannelMessageSend(channel, ErrorMsg)
			return
		}
		getMediaPost(discord, channel, nsfw, subs, "hot")
	case (command == "!joke" || command == "!text") && len(commandContent) == 1:
		subs = SubMap["text"]
		getTextPost(discord, channel, nsfw, subs, "hot")
	case (command == "!joke" || command == "!text") && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		if subs == nil {
			discord.ChannelMessageSend(channel, ErrorMsg)
			return
		}
		getTextPost(discord, channel, nsfw, subs, "hot")
	case (command == "!news" || command == "!link") && len(commandContent) == 1:
		subs = SubMap["news"]
		getLinkPost(discord, channel, nsfw, subs, "hot")
	case (command == "!news" || command == "!link") && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		if subs == nil {
			discord.ChannelMessageSend(channel, ErrorMsg)
			return
		}
		getLinkPost(discord, channel, nsfw, subs, "hot")
	case command == "!fiftyfifty" || command == "!5050":
		subs = []string{"fiftyfifty"}
		getLinkPost(discord, channel, nsfw, subs, "hot")
	case commandContent[0] == "!buzzword":
		getBuzzWord(discord, channel)
	case commandContent[0] == "!hentai":
		// This is still only here because a friend of mine
		// suggested this and I am a nice person
		subs = SubMap["hentai"]
		getMediaPost(discord, channel, nsfw, subs, "hot")
	case command == "!all":
		randchoice := rand.Intn(4)
		switch randchoice {
		case 0:
			getLinkPost(discord, channel, nsfw, []string{"all"}, "")
		case 1:
			getTextPost(discord, channel, nsfw, []string{"all"}, "")
		default:
			getMediaPost(discord, channel, nsfw, []string{"all"}, "")
		}
	case command == "!source":
		err := getSource(discord, channel)
		if err != nil {
			fmt.Println("Error getting source of meme:", err)
			discord.ChannelMessageSend(channel, "Error getting source of meme: "+err.Error())
			return
		}
	case command == "!revsearch":
		imageSearchCommand(discord, channel)
	case command == "!quickmeme":
		var subcommand string
		if len(commandContent) > 1 {
			subcommand = commandContent[1]
		} else {
			subcommand = ""
		}
		subcommand = textFilter(subcommand)
		if !stringInSlice(user.ID, adminIDs) && !isUserMemeBotAdmin(discord, guildID, user) {
			fmt.Println("Intruder tried to execute admin only command:")
			fmt.Println(user.Username)
		} else if stringInSlice(user.ID, adminIDs) {
			switch subcommand {
			case "test":
				quickMemeTest(discord, channel)
			case "testredis":
				quickMemeTestRedis(discord, channel)
			case "testcommoncache":
				quickMemeTestCommonCache(discord, channel)
			case "getcache":
				quickMemeGetCache(discord, channel)
			case "clearcache":
				quickMemeClearCache(discord, channel)
			case "getservercache":
				quickMemeServerCache(discord, channel)
			case "resetblacklist":
				ResetBlacklist()
				discord.ChannelMessageSend(channel, "Blacklist reset. New Blacklist time is "+strconv.FormatInt(BlacklistTime, 10)+".")
			case "ban":
				banSubRoutine(discord, channel, commandContent, guildID, user)
			case "unban":
				unbanSubRoutine(discord, channel, commandContent, guildID, user)
			case "getbanned":
				getbannedSubRoutine(discord, channel, commandContent, guildID, user)
			case "subscribe":
				setQueueRoutine(discord, channel, commandContent, nsfw)
			case "unsubscribe":
				delQueueRoutine(discord, channel)
			default:
				quickMemeDefault(discord, channel)
			}
		} else {
			switch subcommand {
			case "ban":
				banSubRoutine(discord, channel, commandContent, guildID, user)
			case "unban":
				unbanSubRoutine(discord, channel, commandContent, guildID, user)
			case "getbanned":
				getbannedSubRoutine(discord, channel, commandContent, guildID, user)
			case "subscribe":
				setQueueRoutine(discord, channel, commandContent, nsfw)
			case "unsubscribe":
				delQueueRoutine(discord, channel)
			default:
				quickMemeDefault(discord, channel)
			}
		}
	}
	fmt.Println("Posted.")

}
