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
	//LastPost gets the last post from the specified channel string
	LastPost map[string]QuickPost
)

func main() {
	var err error
	var file string
	var key string
	var adminRawIDs []string
	ServerMap = make(map[string]string)
	PostCache = make(map[string][]QuickPost)
	Blacklist = make(map[string][]QuickPost)
	LastPost = make(map[string]QuickPost)
	file = "data.json"
	key, adminRawIDs, err = jsonExtract(file)
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
	commands := []string{"!meme", "!joke", "!hentai", "!news", "!fiftyfifty", "!5050", "!all", "!quickmeme", "!text", "!link", "!source"}
	user := message.Author
	content := message.Content
	guildID := message.GuildID
	if user.ID == botID || user.Bot {
		return
	} else if !ContainsAnySubstring(content, commands) {
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
	commandContent := strings.Split(content, " ")
	sort = "hot"
	switch {
	case commandContent[0] == "!meme" && len(commandContent) == 1:
		subs = []string{"dankmemes", "funny", "memes", "comedyheaven", "MemeEconomy", "therewasanattempt", "wholesomememes", "instant_regret"}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case commandContent[0] == "!meme" && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case (commandContent[0] == "!joke" || commandContent[0] == "!text") && len(commandContent) == 1:
		subs = []string{"jokes", "darkjokes", "antijokes"}
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (commandContent[0] == "!joke" || commandContent[0] == "!text") && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (commandContent[0] == "!news" || commandContent[0] == "!link") && len(commandContent) == 1:
		subs = []string{"UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case (commandContent[0] == "!news" || commandContent[0] == "!link") && len(commandContent) >= 2:
		subs = textFilterSlice(commandContent[1:])
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case commandContent[0] == "!fiftyfifty" || commandContent[0] == "!5050":
		subs = []string{"fiftyfifty"}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case commandContent[0] == "!hentai":
		// This is still only here because a friend of mine
		// suggested this and I am a nice person
		subs = []string{"ahegao", "Artistic_Hentai", "Hentai", "MonsterGirl", "slimegirls", "wholesomehentai", "quick_hentai", "HentaiParadise"}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case commandContent[0] == "!all":
		randchoice := rand.Intn(4)
		switch randchoice {
		case 0:
			err = getLinkPost(discord, channel, nsfw, []string{"all"}, "")
		case 1:
			err = getTextPost(discord, channel, nsfw, []string{"all"}, "")
		default:
			err = getMediaPost(discord, channel, nsfw, []string{"all"}, "")
		}
	case commandContent[0] == "!source":
		err = getSource(discord, channel)
	case commandContent[0] == "!quickmeme":
		var thing string
		if len(commandContent) > 1 {
			thing = commandContent[1]
		} else {
			thing = "status"
		}
		thing = textFilter(thing)
		if !stringInSlice(user.ID, adminIDs) && !isUserMemeBotAdmin(discord, guildID, user) {
			thing = ""
			fmt.Println("Intruder tried to execute admin only command:")
			fmt.Println(user.Username)
		} else if stringInSlice(user.ID, adminIDs) {
			switch thing {
			case "test":
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
			case "getcache":
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
			case "clearcache":
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
			case "getservercache":
				channelCount := len(ServerMap)
				fmt.Println("There are " + strconv.Itoa(channelCount) + " text channels currently cached.")
				fmt.Println(ServerMap)
				discord.ChannelMessageSend(channel, "There are "+strconv.Itoa(channelCount)+" text channels currently cached.")
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
				if len(commandContent) != 3 {
					discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme getbanned [mode]`\nMode can be `channel` or `server`.")
				} else {
					switch commandContent[2] {
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
			switch thing {
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
				if len(commandContent) != 3 {
					discord.ChannelMessageSend(channel, "Incorrect command syntax! Correct syntax is `!quickmeme getbanned [mode]`\nMode can be `channel` or `server`.")
				} else {
					switch commandContent[2] {
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
