package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

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
)

func main() {
	var err error
	var file string
	var key string
	var adminRawIDs []string
	ServerMap = make(map[string]string)
	PostCache = make(map[string][]QuickPost)
	Blacklist = make(map[string][]QuickPost)
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
	commands := []string{"!meme", "!joke", "!hentai", "!news", "!fiftyfifty", "!5050", "!all", "!quickmeme", "!text", "!link"}
	user := message.Author
	content := message.Content
	guildID := message.GuildID
	if user.ID == botID || user.Bot {
		return
	} else if !ContainsAnySubstring(content, commands) {
		return
	}
	channel := message.ChannelID
	if dm {
		channelName = user.Username + "'s DMs"
	} else {
		channelName = "#" + getChannelName(discord, channel, guildID)
	}
	fmt.Println("Command '" + content + "' from " + user.Username + " on " + channelName + " (" + channel + ")")
	nsfw := strings.Contains(strings.ToLower(channelName), "nsfw") || dm
	commandContent := strings.Split(content, " ")
	sort = "hot"
	switch {
	case strings.HasPrefix(content, "!meme") && len(commandContent) == 1:
		subs = []string{"dankmemes", "funny", "memes", "comedyheaven", "MemeEconomy", "therewasanattempt", "wholesomememes", "instant_regret"}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!meme") && len(commandContent) >= 2:
		sub := commandContent[1]
		sub = textFilter(sub)
		subs = []string{sub}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!joke") || strings.HasPrefix(content, "!text")) && len(commandContent) == 1:
		subs = []string{"jokes", "darkjokes", "antijokes"}
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!joke") || strings.HasPrefix(content, "!text")) && len(commandContent) >= 2:
		sub := commandContent[1]
		sub = textFilter(sub)
		subs = []string{sub}
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!news") || strings.HasPrefix(content, "!link")) && len(commandContent) == 1:
		subs = []string{"UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!news") || strings.HasPrefix(content, "!link")) && len(commandContent) >= 2:
		sub := commandContent[1]
		sub = textFilter(sub)
		subs = []string{sub}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!fiftyfifty") || strings.HasPrefix(content, "!5050"):
		subs = []string{"fiftyfifty"}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!hentai"):
		// This is still only here because a friend of mine suggested this
		//fmt.Println("Case 8")
		subs = []string{"ahegao", "Artistic_Hentai", "Hentai", "MonsterGirl", "slimegirls", "wholesomehentai", "quick_hentai", "HentaiParadise"}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!all"):
		randchoice := rand.Intn(4)
		switch randchoice {
		case 0:
			err = getLinkPost(discord, channel, nsfw, []string{"all"}, "")
		case 1:
			err = getTextPost(discord, channel, nsfw, []string{"all"}, "")
		default:
			err = getMediaPost(discord, channel, nsfw, []string{"all"}, "")
		}

	case strings.HasPrefix(content, "!quickmeme"):
		thing := commandContent[1]
		thing = textFilter(thing)
		if !stringInSlice(user.ID, adminIDs) {
			thing = ""
			fmt.Println("Intruder tried to execute admin only command:")
			fmt.Println(user.Username)
		}
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
		default:
			servers := discord.State.Guilds
			userCount := getNumberOfUsers(discord)
			msg := "Discord-Quick-Meme is active and ready on " + strconv.Itoa(len(servers)) + " servers for " + strconv.Itoa(userCount) + " users."
			fmt.Println(msg)
			discord.ChannelMessageSend(channel, msg)
		}
	}
	fmt.Println("Posted.")
	errCheck("Error gettings post info:", err)
}

// Server server object for the golang channels
type Server struct {
	IDs   []string
	Names []string
}

func getAllWorker(discord *discordgo.Session, guildID string, send chan<- Server, wg *sync.WaitGroup, workerNumber int) {
	defer wg.Done()
	var ids []string
	var names []string
	channels, _ := discord.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		ids = append(ids, channel.ID)
		names = append(names, channel.Name)
	}
	server := Server{
		IDs:   ids,
		Names: names,
	}
	send <- server
}

func getAllChannelNames(discord *discordgo.Session) {
	var wg sync.WaitGroup
	fmt.Println("Getting current channel names...")
	starttime := GetMillis()
	guilds := discord.State.Guilds
	bufferSize := len(guilds)
	recv := make(chan Server, bufferSize)
	for i, guild := range guilds {
		wg.Add(1)
		go getAllWorker(discord, guild.ID, recv, &wg, i)
	}
	wg.Wait()
	close(recv)
	for i := 0; i < bufferSize; i++ {
		thing := <-recv
		length := len(thing.IDs)
		for x := 0; x < length; x++ {
			ServerMap[thing.IDs[x]] = thing.Names[x]
		}
	}
	endtime := GetMillis()
	t := endtime - starttime
	fmt.Println("Time to get all current channel names: " + strconv.FormatInt(t, 10) + "ms")
}

func getChannelName(discord *discordgo.Session, channelid string, guildID string) string {
	fmt.Println("Getting channel name....")
	if _, ok := ServerMap[channelid]; ok {
		return ServerMap[channelid]
	}
	starttime := GetMillis()
	channels, _ := discord.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.ID == channelid {
			ServerMap[channelid] = channel.Name
			endtime := GetMillis()
			t := endtime - starttime
			fmt.Println("Time to get channel name: " + strconv.FormatInt(t, 10) + "ms")
			return channel.Name
		}
	}

	return ""
}

func updateStatus(discord *discordgo.Session) {
	uCount := getNumberOfUsers(discord)
	err := discord.UpdateStatus(0, "with "+strconv.Itoa(uCount)+" others")
	if err != nil {
		errCheck("Error attempting to set the status.", err)
	}
}

func stringInSlice(s string, a []string) bool {
	for _, thing := range a {
		if thing == s {
			return true
		}
	}
	return false
}

// ComesFromDM returns true if a message comes from a DM channel
func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func textFilter(input string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	outputString := reg.ReplaceAllString(input, "")
	return outputString
}

func getNumberOfUsers(discord *discordgo.Session) int {
	count := 0
	for _, guild := range discord.State.Guilds {
		count += len(guild.Members)
	}
	return count
}

// GetMillis gets number of milliseconds since epoch as a 64bit integer
func GetMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	return nanos / 1000000
}

// ContainsAnySubstring Checks if any of the strings in the array are in the test string
func ContainsAnySubstring(testString string, strArray []string) bool {
	for _, str := range strArray {
		if strings.Contains(testString, str) {
			return true
		}
	}
	return false
}

func getMediaPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) error {
	var returnPost QuickPost
	var err error
	var score int32
	var url string
	var title string
	var nsfw bool
	var postlink string
	var sub string
	rand.Seed(time.Now().Unix())
	randColor := rand.Intn(0xffffff)
	imageEndings := []string{".jpg", ".png", ".jpeg"}
	limit := 100
	toggled := false
	for i := 0; i < 5; i++ {
		returnPost, sub = GetPost(subs, limit, sort, "media")
		blacklisted := CheckBlacklist(channel, returnPost)
		score = returnPost.Score
		url = returnPost.Content
		title = returnPost.Title
		postlink = returnPost.Permalink
		nsfw = returnPost.Nsfw
		if channelNsfw && !blacklisted {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if channelNsfw && !nsfw && !blacklisted {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if !channelNsfw && !nsfw && !blacklisted {
			toggled = true

			break
		} else {
			if !blacklisted {
				fmt.Println("Channel is not NSFW but post is NSFW, retrying...")
			}
		}
	}
	if ContainsAnySubstring(url, imageEndings) && toggled {
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       randColor,
			Description: "Score: " + strconv.FormatInt(int64(score), 10),
			Image: &discordgo.MessageEmbedImage{
				URL: url,
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     title,
		}
		_, err = discord.ChannelMessageSend(channel, "From r/"+sub)
		_, err = discord.ChannelMessageSendEmbed(channel, embed)
		_, err = discord.ChannelMessageSend(channel, "From https://reddit.com"+postlink)
	} else if toggled {
		_, err = discord.ChannelMessageSend(channel, "From r/"+sub)
		_, err = discord.ChannelMessageSend(channel, url)
		_, err = discord.ChannelMessageSend(channel, title)
		_, err = discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10)+"\nOriginal Post: https://reddit.com"+postlink)
	} else {
		_, err = discord.ChannelMessageSend(channel, "Error!")
		_, err = discord.ChannelMessageSend(channel, "Too many tries to not find NSFW post, maybe that Subreddit is filled with them? Hint: Add \"NSFW\" to the channel name to allow NSFW posts.")
	}

	return err
}

func getTextPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) error {
	var returnPost QuickPost
	var err error
	var score int32
	var text string
	var title string
	var nsfw bool
	var postlink string
	var sub string
	limit := 100
	toggled := false
	for i := 0; i < 10; i++ {
		returnPost, sub = GetPost(subs, limit, sort, "text")
		blacklisted := CheckBlacklist(channel, returnPost)
		score = returnPost.Score
		text = returnPost.Content
		title = returnPost.Title
		postlink = returnPost.Permalink
		nsfw = returnPost.Nsfw
		if channelNsfw && !blacklisted {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if channelNsfw && !nsfw && !blacklisted {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if !channelNsfw && !nsfw && !blacklisted {
			toggled = true

			break
		} else {
			if !blacklisted {
				fmt.Println("Channel is not NSFW but post is NSFW, retrying...")
			}
		}
	}
	if toggled {
		_, err = discord.ChannelMessageSend(channel, "From r/"+sub)
		_, err = discord.ChannelMessageSend(channel, title)
		_, err = discord.ChannelMessageSend(channel, text)
		_, err = discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10)+"\nOriginal Post: https://reddit.com"+postlink)
	} else {
		_, err = discord.ChannelMessageSend(channel, "Error!")
		_, err = discord.ChannelMessageSend(channel, "Too many tries to not find NSFW post, maybe that Subreddit is filled with them? Hint: Add \"NSFW\" to the channel name to allow NSFW posts.")
	}
	return err
}

func getLinkPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string, sort string) error {
	var returnPost QuickPost
	var err error
	var score int32
	var url string
	var title string
	var nsfw bool
	var postlink string
	var sub string
	limit := 100
	toggled := false
	for i := 0; i < 10; i++ {
		returnPost, sub = GetPost(subs, limit, sort, "link")
		blacklisted := CheckBlacklist(channel, returnPost)
		score = returnPost.Score
		url = returnPost.Content
		title = returnPost.Title
		postlink = returnPost.Permalink
		nsfw = returnPost.Nsfw
		if channelNsfw && !blacklisted {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if channelNsfw && !nsfw && !blacklisted {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if !channelNsfw && !nsfw && !blacklisted {
			toggled = true

			break
		} else {
			if !blacklisted {
				fmt.Println("Channel is not NSFW but post is NSFW, retrying...")
			}
		}
	}

	if toggled {
		_, err = discord.ChannelMessageSend(channel, "From r/"+sub)
		_, err = discord.ChannelMessageSend(channel, url)
		_, err = discord.ChannelMessageSend(channel, title)
		_, err = discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10)+"\nOriginal Post: https://reddit.com"+postlink)
	} else {
		_, err = discord.ChannelMessageSend(channel, "Error!")
		_, err = discord.ChannelMessageSend(channel, "Too many tries to not find NSFW post, maybe that Subreddit is filled with them? Hint: Add \"NSFW\" to the channel name to allow NSFW posts.")
	}
	return err
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
