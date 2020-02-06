package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	botID         string
	adminIDs      []string
	// CacheTime stores cache timer value
	CacheTime int64
	//BlacklistTime stores the blacklist time for all of the channels
	BlacklistTime int64
	// ServerMap this is all of the servers an the servers this gets wiped from memory as soon as the Bot gets killed
	ServerMap map[string]string
	//PostCache stores all posts
	PostCache map[string][]QuickPost
	//Blacklist list of all of the post that are blacklisted from the specified channel
	Blacklist map[string][]QuickPost // will be wiped every two to three hours
	//CommonSubs stores the amount of times the subs are hit
	CommonSubs map[string]uint8 // only needs to count up to 10
	//CommonSubsTime if one week passes, clear the above cache
	// 604800000ms in a week
	CommonSubsTime    map[string]int64
	CommonSubsCounter uint8
	//LastPost gets the last post from the specified channel string
	LastPost map[string]QuickPost
	//SubMap contains all of the data for the subs
	SubMap map[string][]string
	//CachePopulating if true, do not run the populate cache until finished
	CachePopulating bool
	mrisaAddress    string
)

func main() {
	var err error
	var file string
	var key string
	var adminRawIDs []string
	ServerMap = make(map[string]string)
	PostCache = make(map[string][]QuickPost)
	Blacklist = make(map[string][]QuickPost)
	CommonSubs = make(map[string]uint8)
	CommonSubsTime = make(map[string]int64)
	CommonSubsCounter = 0
	LastPost = make(map[string]QuickPost)
	SubMap = make(map[string][]string)
	file = "data.json"
	key, adminRawIDs, err = loginExtract(file)
	mrisaAddress = mrisaExtract(file)
	errCheck("Error opening key file", err)
	discord, err := discordgo.New("Bot " + key)
	errCheck("Error creating discord session", err)
	user, err := discord.User("@me")
	for _, admin := range adminRawIDs {
		a, _ := discord.User(admin)
		adminIDs = append(adminIDs, a.ID)
	}
	errCheck("error retrieving account", err)
	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(readyHandler)
	err = discord.Open()
	errCheck("Error opening discord connection", err)
	defer discord.Close()
	commandPrefix = "!"
	<-make(chan struct{})
}

func readyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	servers := discord.State.Guilds
	getAllChannelNames(discord)
	CachePopulating = true
	PopulateCache()
	ResetBlacklist()
	fmt.Println("Discord-Quick-Meme has started on " + strconv.Itoa(len(servers)) + " servers")
	go updateStatus(discord)
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	var err error
	var sort string
	var subs []string
	var dm bool
	var channelName string
	go updateStatus(discord)
	go UpdateBlacklistTime()
	dm, err = ComesFromDM(discord, message)
	commands := []string{"!meme", "!joke", "!hentai", "!news", "!fiftyfifty", "!5050", "!all", "!quickmeme", "!text", "!link", "!source", "!buzzword", "!search"}
	user := message.Author
	content := message.Content
	commandContent := strings.Split(content, " ")
	command := commandContent[0]
	guildID := message.GuildID
	if user.ID == botID || user.Bot || !stringInSlice(commandContent[0], commands) {
		return
	}
	channelObject, _ := discord.Channel(message.ChannelID)
	channel := message.ChannelID
	if dm {
		channelName = user.Username + "'s DMs"
	} else {
		channelName = "#" + getChannelName(discord, channel, guildID)
	}
	fmt.Println("Command '" + content + "' from " + user.Username + " on " + channelName + " (" + channel + ")")
	nsfw := channelObject.NSFW || dm
	sort = "hot"
	switch {
	case command == "!meme" && len(commandContent) == 1:
		subs = SubMap["memes"]
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case command == "!meme" && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case (command == "!joke" || command == "!text") && len(commandContent) == 1:
		subs = SubMap["text"]
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (command == "!joke" || command == "!text") && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (command == "!news" || command == "!link") && len(commandContent) == 1:
		subs = SubMap["news"]
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case (command == "!news" || command == "!link") && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case command == "!fiftyfifty" || command == "!5050":
		subs = []string{"fiftyfifty"}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case commandContent[0] == "!buzzword":
		err = getBuzzWord(discord, channel)
	case commandContent[0] == "!hentai":
		// This is still only here because a friend of mine
		// suggested this and I am a nice person
		subs = SubMap["hentai"]
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case command == "!all":
		randchoice := rand.Intn(4)
		switch randchoice {
		case 0:
			err = getLinkPost(discord, channel, nsfw, []string{"all"}, "")
		case 1:
			err = getTextPost(discord, channel, nsfw, []string{"all"}, "")
		default:
			err = getMediaPost(discord, channel, nsfw, []string{"all"}, "")
		}
	case command == "!source":
		err = getSource(discord, channel)
	case command == "!search":
		quickMemeImageSearch(discord, channel)
	case command == "!quickmeme":
		var subcommand string
		if len(commandContent) > 1 {
			subcommand = commandContent[1]
		} else {
			subcommand = "status"
		}
		subcommand = textFilter(subcommand)
		if !stringInSlice(user.ID, adminIDs) && !isUserMemeBotAdmin(discord, guildID, user) {
			subcommand = ""
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
				if len(commandContent) < 4 || len(commandContent) > 5 {
					discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme ban [mode] [subreddit]`\nMode can be `channel` or `server`.")
				} else if isUserMemeBotAdmin(discord, guildID, user) { // fix this
					switch commandContent[2] {
					case "server":
						channels, _ := discord.GuildChannels(guildID)
						subreddits := textFilterSlice(commandContent[3:])
						for _, chat := range channels {
							// this should be async to save time
							for _, subreddit := range subreddits {
								go AppendBannedSubreddit(chat.ID, subreddit)
							}
						}
						// this should be a message about the ban
						discord.ChannelMessageSend(channel, user.Mention()+" banned subreddit(s) "+strings.Join(subreddits, ", ")+" on all channels.")
					default:
						subreddits := textFilterSlice(commandContent[3:])
						for _, subreddit := range subreddits {
							go AppendBannedSubreddit(channel, subreddit)
						}
						discord.ChannelMessageSend(channel, user.Mention()+" banned subreddit(s) "+strings.Join(subreddits, ", ")+".")
					}
				} else {
					discord.ChannelMessageSend(channel, "Insufficient Permissions! You must have the \"Memebot Admin\" role to ban subreddits!")
				}
			case "unban":
				if len(commandContent) < 4 || len(commandContent) > 5 {
					discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme unban [mode] [subreddit]`\nMode can be `channel` or `server`.")
				} else if isUserMemeBotAdmin(discord, guildID, user) { // fix this
					switch commandContent[2] {
					case "server":
						channels, _ := discord.GuildChannels(guildID)
						subreddits := textFilterSlice(commandContent[3:])
						for _, chat := range channels {
							// this should be async to save time
							for _, subreddit := range subreddits {
								go UnbanSubreddit(chat.ID, subreddit)
							}
						}
						// this should be a message about the ban
						discord.ChannelMessageSend(channel, user.Mention()+" unbanned subreddit(s) "+strings.Join(subreddits, ", ")+" on all channels.")
					default:
						// there should be a message about the ban here
						subreddits := textFilterSlice(commandContent[3:])
						for _, subreddit := range subreddits {
							go UnbanSubreddit(channel, subreddit)
						}
						discord.ChannelMessageSend(channel, user.Mention()+" unbanned subreddit(s) "+strings.Join(subreddits, ", ")+".")
					}
				} else {
					discord.ChannelMessageSend(channel, "Insufficient Permissions! You must have the \"Memebot Admin\" role to ban subreddits!")
				}
			case "getbanned":
				banContext := commandContent[2]
				if len(commandContent) != 3 {
					banContext = "channel"
				} else {
					switch banContext {
					case "server":
						channels, _ := discord.GuildChannels(guildID)
						for _, chat := range channels {
							bannedSubs, err := GetBannedSubreddits(chat.ID)
							if err != nil {
								discord.ChannelMessageSend(channel, "There was an error processing your request. Please report this at https://github.com/chand1012/Discord-Quick-Meme/issues")
								fmt.Println(err)
								break
							}
							msgString := strings.Join(bannedSubs, ", ")
							if msgString != "" && chat.Type == discordgo.ChannelTypeGuildText {
								discord.ChannelMessageSend(channel, "Subs banned on "+chat.Name+":\n"+msgString)
							}
						}
					default:
						bannedSubs, err := GetBannedSubreddits(channel)
						if err != nil {
							discord.ChannelMessageSend(channel, "There was an error processing your request. Please report this at https://github.com/chand1012/Discord-Quick-Meme/issues")
							fmt.Println(err)
						} else {
							msgString := strings.Join(bannedSubs, ", ")
							discord.ChannelMessageSend(channel, "Subs banned on this channel:\n"+msgString)
						}
					}
				}
			default:
				servers := discord.State.Guilds
				userCount := getNumberOfUsers(discord)
				msg := "Discord-Quick-Meme is active and ready on " + strconv.Itoa(len(servers)) + " servers for " + strconv.Itoa(userCount) + " users."
				fmt.Println(msg)
				discord.ChannelMessageSend(channel, msg)
			}
		} else {
			switch subcommand {
			case "ban":
				if len(commandContent) <= 4 {
					discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme ban [mode] [subreddit]`\nMode can be `channel` or `server`.")
				} else if isUserMemeBotAdmin(discord, guildID, user) { // fix this
					switch commandContent[2] {
					case "server":
						channels, _ := discord.GuildChannels(guildID)
						subreddits := textFilterSlice(commandContent[3:])
						for _, chat := range channels {
							// this should be async to save time
							for _, subreddit := range subreddits {
								go AppendBannedSubreddit(chat.ID, subreddit)
							}
						}
						// this should be a message about the ban
						discord.ChannelMessageSend(channel, user.Mention()+" banned subreddit(s) "+strings.Join(subreddits, ", ")+" on all channels.")
					default:
						subreddits := textFilterSlice(commandContent[3:])
						for _, subreddit := range subreddits {
							go AppendBannedSubreddit(channel, subreddit)
						}
						discord.ChannelMessageSend(channel, user.Mention()+" banned subreddit(s) "+strings.Join(subreddits, ", ")+".")
					}
				} else {
					discord.ChannelMessageSend(channel, "Insufficient Permissions! You must have the \"Memebot Admin\" role to ban subreddits!")
				}
			case "unban":
				if len(commandContent) <= 4 {
					discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme unban [mode] [subreddit]`\nMode can be `channel` or `server`.")
				} else if isUserMemeBotAdmin(discord, guildID, user) { // fix this
					switch commandContent[2] {
					case "server":
						channels, _ := discord.GuildChannels(guildID)
						subreddits := textFilterSlice(commandContent[3:])
						for _, chat := range channels {
							// this should be async to save time
							for _, subreddit := range subreddits {
								go UnbanSubreddit(chat.ID, subreddit)
							}
						}
						// this should be a message about the ban
						discord.ChannelMessageSend(channel, user.Mention()+" unbanned subreddit(s) "+strings.Join(subreddits, ", ")+" on all channels.")
					default:
						// there should be a message about the ban here
						subreddits := textFilterSlice(commandContent[3:])
						for _, subreddit := range subreddits {
							go UnbanSubreddit(channel, subreddit)
						}
						discord.ChannelMessageSend(channel, user.Mention()+" unbanned subreddit(s) "+strings.Join(subreddits, ", ")+".")
					}
				} else {
					discord.ChannelMessageSend(channel, "Insufficient Permissions! You must have the \"Memebot Admin\" role to ban subreddits!")
				}
			case "getbanned":
				banContext := commandContent[2]
				if len(commandContent) != 3 {
					banContext = "channel"
				} else {
					switch banContext {
					case "server":
						channels, _ := discord.GuildChannels(guildID)
						for _, chat := range channels {
							bannedSubs, err := GetBannedSubreddits(chat.ID)
							if err != nil {
								discord.ChannelMessageSend(channel, "There was an error processing your request. Please report this at https://github.com/chand1012/Discord-Quick-Meme/issues")
								fmt.Println(err)
								break
							}
							msgString := strings.Join(bannedSubs, ", ")
							if msgString != "" && chat.Type == discordgo.ChannelTypeGuildText {
								discord.ChannelMessageSend(channel, "Subs banned on "+chat.Name+":\n"+msgString)
							}
						}
					default:
						bannedSubs, err := GetBannedSubreddits(channel)
						if err != nil {
							discord.ChannelMessageSend(channel, "There was an error processing your request. Please report this at https://github.com/chand1012/Discord-Quick-Meme/issues")
							fmt.Println(err)
						} else {
							msgString := strings.Join(bannedSubs, ", ")
							discord.ChannelMessageSend(channel, "Subs banned on this channel:\n"+msgString)
						}
					}
				}
			default:
				servers := discord.State.Guilds
				userCount := getNumberOfUsers(discord)
				msg := "Discord-Quick-Meme is active and ready on " + strconv.Itoa(len(servers)) + " servers for " + strconv.Itoa(userCount) + " users."
				fmt.Println(msg)
				discord.ChannelMessageSend(channel, msg)
			}
		}
	}
	fmt.Println("Posted.")
	errCheck("Error gettings post info:", err)
}
