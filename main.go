package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	botID         string
	// ServerMap
	// this is all of the servers an the servers
	// this gets wiped from memory as soon as the Bot gets killed
	ServerMap map[string]string
)

func main() {
	var err error
	var file string
	var key string
	ServerMap = make(map[string]string)
	file = "data.json"
	key, _, err = jsonExtract(file)
	errCheck("Error opening key file", err)
	discord, err := discordgo.New("Bot " + key)
	errCheck("Error creating discord session", err)
	user, err := discord.User("@me")
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
	err := discord.UpdateStatus(0, "with spacetime.")
	if err != nil {
		fmt.Println("Error attempting to set the status.")
	}
	servers := discord.State.Guilds
	fmt.Println("Discord-Quick-Meme has started on " + strconv.Itoa(len(servers)) + " servers")
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	var err error
	var subs []string
	commands := []string{"!meme", "!joke", "!hentai", "!news", "!fiftyfifty", "!5050", "!all"}
	user := message.Author
	content := message.Content
	if user.ID == botID || user.Bot {
		return
	} else if !ContainsAnySubstring(content, commands) {
		return
	}
	channel := message.ChannelID
	channelName := getChannelName(discord, channel)
	fmt.Println("Command '" + content + "' from #" + channelName + " (" + channel + ")")
	nsfw := strings.Contains(strings.ToLower(channelName), "nsfw")
	contentLength := utf8.RuneCountInString(content)
	if strings.HasPrefix(content, "!meme") && contentLength <= 5 {
		subs = []string{"dankmemes", "funny", "memes", "dank_meme", "comedyheaven", "CyanideandHappiness", "therewasanattempt", "wholesomememes", "instant_regret"}
		err = getMediaPost(discord, channel, nsfw, subs)
	} else if strings.HasPrefix(content, "!meme") && contentLength >= 5 {
		sub := content[6:]
		subs = []string{sub}
		err = getMediaPost(discord, channel, nsfw, subs)
	} else if strings.HasPrefix(content, "!joke") && contentLength <= 5 {
		subs = []string{"jokes", "darkjokes", "antijokes"}
		err = getTextPost(discord, channel, nsfw, subs)
	} else if (strings.HasPrefix(content, "!joke") || strings.HasPrefix(content, "!text")) && contentLength >= 5 {
		sub := content[6:]
		subs = []string{sub}
		err = getTextPost(discord, channel, nsfw, subs)
	} else if strings.HasPrefix(content, "!news") {
		subs = []string{"UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"}
		err = getLinkPost(discord, channel, nsfw, subs)
	} else if strings.HasPrefix(content, "!fiftyfifty") || strings.HasPrefix(content, "!5050") {
		subs = []string{"fiftyfifty"}
		err = getLinkPost(discord, channel, nsfw, subs)
	} else if strings.HasPrefix(content, "!hentai") {
		// This is still only here because a friend of mine suggested this
		subs = []string{"ahegao", "Artistic_Hentai", "Hentai", "MonsterGirl", "slimegirls", "wholesomehentai", "quick_hentai", "HentaiParadise"}
		err = getMediaPost(discord, channel, nsfw, subs)
	} else if strings.HasPrefix(content, "!all") {
		randchoice := rand.Intn(4)
		if randchoice == 0 {
			err = getLinkPost(discord, channel, nsfw, []string{"all"})
		} else if randchoice == 1 {
			err = getTextPost(discord, channel, nsfw, []string{"all"})
		} else {
			err = getMediaPost(discord, channel, nsfw, []string{"all"})
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
		for _, guild := range discord.State.Guilds {
			channels, _ := discord.GuildChannels(guild.ID)

			for _, channel := range channels {
				if channel.ID == channelid {
					ServerMap[channelid] = channel.Name
					return channel.Name
				}
			}
		}
	}
	return ""
}

// ContainsAnySubstring checks if any of the strings in the array are in the test string
func ContainsAnySubstring(testString string, strArray []string) bool {
	for _, str := range strArray {
		if strings.Contains(testString, str) {
			return true
		}
	}
	return false
}

func getMediaPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string) error {
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
	limit := 25
	toggled := false
	for i := 0; i < 5; i++ {
		score, url, title, nsfw, postlink, sub = GetMediaPost(subs, limit)
		if channelNsfw {
			toggled = true
			break
		} else if channelNsfw && !nsfw {
			toggled = true
			break
		} else if !channelNsfw && !nsfw {
			toggled = true
			break
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

func getTextPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string) error {
	var err error
	var score int32
	var text string
	var title string
	var nsfw bool
	var postlink string
	var sub string
	limit := 25
	toggled := false
	for i := 0; i < 10; i++ {
		score, text, title, nsfw, postlink, sub = GetTextPost(subs, limit)
		if channelNsfw {
			toggled = true
			break
		} else if channelNsfw && !nsfw {
			toggled = true
			break
		} else if !channelNsfw && !nsfw {
			toggled = true
			break
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

func getLinkPost(discord *discordgo.Session, channel string, channelNsfw bool, subs []string) error {
	var err error
	var score int32
	var url string
	var title string
	var nsfw bool
	var postlink string
	var sub string
	limit := 25
	toggled := false
	for i := 0; i < 10; i++ {
		score, url, title, nsfw, postlink, sub = GetMediaPost(subs, limit)
		if channelNsfw {
			toggled = true
			break
		} else if channelNsfw && !nsfw {
			toggled = true
			break
		} else if !channelNsfw && !nsfw {
			toggled = true
			break
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
