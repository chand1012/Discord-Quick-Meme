package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	botID         string
	adminID       string
	// CacheTime stores cache timer value
	CacheTime int64
	// ServerMap this is all of the servers an the servers this gets wiped from memory as soon as the Bot gets killed
	ServerMap map[string]string
	//PostCache stores all posts
	PostCache map[string][]QuickPost
)

func main() {
	var err error
	var file string
	var key string
	var adminUsername string
	ServerMap = make(map[string]string)
	PostCache = make(map[string][]QuickPost)
	file = "data.json"
	key, adminUsername, err = jsonExtract(file)
	errCheck("Error opening key file", err)
	discord, err := discordgo.New("Bot " + key)
	errCheck("Error creating discord session", err)
	user, err := discord.User("@me")
	admin, err := discord.User(adminUsername)
	errCheck("error retrieving account", err)
	botID = user.ID
	adminID = admin.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(readyHandler)
	err = discord.Open()
	errCheck("Error opening discord connection", err)
	defer discord.Close()
	commandPrefix = "!"
	<-make(chan struct{})

}

func readyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	err := discord.UpdateStatus(0, "with spacetime.")
	if err != nil {
		fmt.Println("Error attempting to set the status.")
	}
	servers := discord.State.Guilds
	fmt.Println("Discord-Quick-Meme has started on " + strconv.Itoa(len(servers)) + " servers")
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	var err error
	var sort string
	var subs []string
	var dm bool
	var channelName string
	dm, err = ComesFromDM(discord, message)
	commands := []string{"!meme", "!joke", "!hentai", "!news", "!fiftyfifty", "!5050", "!all", "!quickmeme", "!text", "!link"}
	user := message.Author
	content := message.Content
	if user.ID == botID || user.Bot {
		return
	} else if !ContainsAnySubstring(content, commands) {
		return
	}
	channel := message.ChannelID
	if dm {
		channelName = user.Username + "'s DMs"
	} else {
		channelName = "#" + getChannelName(discord, channel)
	}
	fmt.Println("Command '" + content + "' from " + user.Username + " on " + channelName + " (" + channel + ")")
	nsfw := strings.Contains(strings.ToLower(channelName), "nsfw")
	contentLength := utf8.RuneCountInString(content)
	sort = "hot"
	switch {
	case strings.HasPrefix(content, "!meme") && contentLength <= 5:
		//fmt.Println("Case 1")
		subs = []string{"dankmemes", "funny", "memes", "comedyheaven", "blackpeopletwitter", "whitepeopletwitter", "MemeEconomy", "therewasanattempt", "wholesomememes", "instant_regret"}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!meme") && contentLength >= 5:
		//fmt.Println("Case 2")
		sub := content[6:]
		sub, err = textFilter(sub)
		subs = []string{sub}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!joke") || strings.HasPrefix(content, "!text")) && contentLength <= 5:
		//fmt.Println("Case 3")
		subs = []string{"jokes", "darkjokes", "antijokes"}
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!joke") || strings.HasPrefix(content, "!text")) && contentLength >= 5:
		//fmt.Println("Case 4")
		sub := content[6:]
		sub, err = textFilter(sub)
		subs = []string{sub}
		err = getTextPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!news") || strings.HasPrefix(content, "!link")) && contentLength <= 5:
		//fmt.Println("Case 5")
		subs = []string{"UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case (strings.HasPrefix(content, "!news") || strings.HasPrefix(content, "!link")) && contentLength >= 5:
		//fmt.Println("Case 6")
		sub := content[6:]
		sub, err = textFilter(sub)
		subs = []string{sub}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!fiftyfifty") || strings.HasPrefix(content, "!5050"):
		//fmt.Println("Case 7")
		subs = []string{"fiftyfifty"}
		err = getLinkPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!hentai"):
		// This is still only here because a friend of mine suggested this
		//fmt.Println("Case 8")
		subs = []string{"ahegao", "Artistic_Hentai", "Hentai", "MonsterGirl", "slimegirls", "wholesomehentai", "quick_hentai", "HentaiParadise"}
		err = getMediaPost(discord, channel, nsfw, subs, sort)
	case strings.HasPrefix(content, "!all"):
		//fmt.Println("Case 9")
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
		//fmt.Println(content)
		thing := content[10:]
		thing, err = textFilter(thing)
		if user.ID != adminID {
			thing = ""
		}
		switch thing {
		default:
			servers := discord.State.Guilds
			userCount := getNumberOfUsers(discord)
			msg := "Discord-Quick-Meme is active and ready on " + strconv.Itoa(len(servers)) + " servers for " + strconv.Itoa(userCount) + " users."
			fmt.Println(msg)
			discord.ChannelMessageSend(channel, msg)
		case " test":
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
		case " getcache":
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
		case " clearcache":
			discord.ChannelMessageSend(channel, "Clearing cache...")
			fmt.Println("Admin issued cache clear...")
			ClearCache()
			CacheTime = time.Now().Unix() + 3600
			msg := "New cache time is " + strconv.FormatInt(CacheTime, 10)
			fmt.Println(msg)
			discord.ChannelMessageSend(channel, "Done. "+msg)
		}
	}
	fmt.Println("Posted.")
	errCheck("Error gettings post info:", err)
}

func getChannelName(discord *discordgo.Session, channelid string) string {
	fmt.Println("Getting channel name....")
	if _, ok := ServerMap[channelid]; ok {
		return ServerMap[channelid]
	} else {
		starttime := GetMillis()
		for _, guild := range discord.State.Guilds {
			channels, _ := discord.GuildChannels(guild.ID)

			for _, channel := range channels {
				if channel.ID == channelid {
					endtime := GetMillis()
					t := endtime - starttime
					fmt.Println("Time took to find channel name: " + strconv.FormatInt(t, 10) + "ms")
					ServerMap[channelid] = channel.Name

					return channel.Name
				}
			}
		}
	}
	return ""
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

func textFilter(input string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	outputString := reg.ReplaceAllString(input, "")
	return outputString, err
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
		returnPost, sub = GetMediaPost(subs, limit, sort)
		score = returnPost.Score
		url = returnPost.Content
		title = returnPost.Title
		postlink = returnPost.Permalink
		nsfw = returnPost.Nsfw
		if channelNsfw {
			toggled = true
			break
		} else if channelNsfw && !nsfw {
			toggled = true
			break
		} else if !channelNsfw && !nsfw {
			toggled = true
			break
		} else {
			fmt.Println("Channel is not NSFW but post is NSFW, retrying...")
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
		_, err = discord.ChannelMessageSend(channel, "Too many tries to not find NSFW post, maybe that Subreddit is filled with them?")
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
		returnPost, sub = GetTextPost(subs, limit, sort)
		score = returnPost.Score
		text = returnPost.Content
		title = returnPost.Title
		postlink = returnPost.Permalink
		nsfw = returnPost.Nsfw
		if channelNsfw {
			toggled = true
			break
		} else if channelNsfw && !nsfw {
			toggled = true
			break
		} else if !channelNsfw && !nsfw {
			toggled = true
			break
		} else {
			fmt.Println("Channel is not NSFW but post is NSFW, retrying...")
		}
	}

	if toggled {
		_, err = discord.ChannelMessageSend(channel, "From r/"+sub)
		_, err = discord.ChannelMessageSend(channel, title)
		_, err = discord.ChannelMessageSend(channel, text)
		_, err = discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10)+"\nOriginal Post: https://reddit.com"+postlink)
	} else {
		_, err = discord.ChannelMessageSend(channel, "Error!")
		_, err = discord.ChannelMessageSend(channel, "Too many tries to not find NSFW post, maybe that Subreddit is filled with them?")
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
		returnPost, sub = GetLinkPost(subs, limit, sort)
		score = returnPost.Score
		url = returnPost.Content
		title = returnPost.Title
		postlink = returnPost.Permalink
		nsfw = returnPost.Nsfw
		if channelNsfw {
			toggled = true
			break
		} else if channelNsfw && !nsfw {
			toggled = true
			break
		} else if !channelNsfw && !nsfw {
			toggled = true
			break
		} else {
			fmt.Println("Channel is not NSFW but post is NSFW, retrying...")
		}
	}

	if toggled {
		_, err = discord.ChannelMessageSend(channel, "From r/"+sub)
		_, err = discord.ChannelMessageSend(channel, url)
		_, err = discord.ChannelMessageSend(channel, title)
		_, err = discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10)+"\nOriginal Post: https://reddit.com"+postlink)
	} else {
		_, err = discord.ChannelMessageSend(channel, "Error!")
		_, err = discord.ChannelMessageSend(channel, "Too many tries to not find NSFW post, maybe that Subreddit is filled with them?")
	}
	return err
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
