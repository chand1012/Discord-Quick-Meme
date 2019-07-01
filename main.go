package main

import (
	"errors"
	"fmt"
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
	fmt.Println("Discord-Quick-Meme has started on " + string(len(servers)) + " servers")
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		return
	}
	// this will be epanded upon in the near future
	content := message.Content
	channel := message.ChannelID
	channelName := getChannelName(discord, channel)
	nsfw := strings.Contains(channelName, "nsfw")
	err := getMemePost(discord, channel, nsfw)
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
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
	subs := []string{"dankmemes", "funny", "memes", "dank_meme", "comedyheaven", "CyanideandHappiness", "therewasanattempt", "wholesomememes", "instant_regret"}
	imageEndings := []string{".jpg", ".png", ".jpeg", ".gif", ".gifv", "gfycat"}
	limit := 100
	toggled := false
	for i := 0; i < 10; i++ {
		score, url, title, nsfw, postlink, sub = GetPost(subs, limit)
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
			Color:       0x00ff00,
			Description: "Score: " + string(score),
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
		discord.ChannelMessageSend(channel, "Score: "+string(score)+"\nOriginal Post: https://reddit.com"+postlink)
	}

	return err
}

/*
func getJokePost(discord *discordgo.Session, channel string) error {

}
*/
func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
