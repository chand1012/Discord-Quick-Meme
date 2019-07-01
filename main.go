package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	botID         string
)

func main() {
	var err error
	var file string
	var key string
	file = "data.json"
	key, err = jsonExtract(file)
	errCheck("Error opening key file", err)
	//fmt.Println(keys.botId)
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
	user := message.Author
	if user.ID == botID || user.Bot {
		return
	}

	content := message.Content
	channel := message.ChannelID
	channelName := getChannelName(discord, channel)
	nsfw := strings.Contains(strings.ToLower(channelName), "nsfw")
	if strings.HasPrefix(content, "!meme") {
		err = getMemePost(discord, channel, nsfw)
	} else if strings.HasPrefix(content, "!joke") {
		err = getJokePost(discord, channel, nsfw)
	}
	errCheck("Error gettings post info:", err)
}

func getChannelName(discord *discordgo.Session, channelid string) string {
	for _, guild := range discord.State.Guilds {
		channels, _ := discord.GuildChannels(guild.ID)

		for _, channel := range channels {
			if channel.ID == channelid {
				return channel.Name
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

func getMemePost(discord *discordgo.Session, channel string, channelNsfw bool) error {
	var err error
	var score int32
	var url string
	var title string
	var nsfw bool
	var postlink string
	var sub string
	rand.Seed(time.Now().Unix())
	randColor := rand.Intn(0xffffff)
	subs := []string{"dankmemes", "funny", "memes", "dank_meme", "comedyheaven", "CyanideandHappiness", "therewasanattempt", "wholesomememes", "instant_regret"}
	imageEndings := []string{".jpg", ".png", ".jpeg"}
	limit := 100
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
	if !toggled {
		err = errors.New("too many attempts to find non nsfw post")
	}
	if ContainsAnySubstring(url, imageEndings) {
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       randColor,
			Description: "Score: " + strconv.FormatInt(int64(score), 10),
			Image: &discordgo.MessageEmbedImage{
				URL: url,
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: url,
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     title,
		}
		discord.ChannelMessageSend(channel, "From r/"+sub)
		discord.ChannelMessageSendEmbed(channel, embed)
		discord.ChannelMessageSend(channel, "From https://reddit.com"+postlink)
	} else {
		discord.ChannelMessageSend(channel, "From r/"+sub)
		discord.ChannelMessageSend(channel, url)
		discord.ChannelMessageSend(channel, title)
		discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10)+"\nOriginal Post: https://reddit.com"+postlink)
	}

	return err
}

func getJokePost(discord *discordgo.Session, channel string, channelNsfw bool) error {
	var err error
	var score int32
	var text string
	var title string
	var nsfw bool
	var postlink string
	var sub string
	subs := []string{"jokes", "darkjokes", "antijokes"}
	limit := 100
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
	if !toggled {
		err = errors.New("too many attempts to find non nsfw post")
	}
	discord.ChannelMessageSend(channel, "From r/"+sub)
	discord.ChannelMessageSend(channel, text)
	discord.ChannelMessageSend(channel, title)
	discord.ChannelMessageSend(channel, "Score: "+strconv.FormatInt(int64(score), 10)+"\nOriginal Post: https://reddit.com"+postlink)
	return err
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
