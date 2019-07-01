package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	botId         string
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
	botId = user.ID
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
	fmt.Printf("Discord-Quick-Meme has started on %d servers", len(servers))
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.Bot {
		return
	}
	// this will be epanded upon in the near future
	content := message.Content
	fmt.Println(content)
}

func getDiscordPost() (*discordgo.MessageEmbed, bool) {
	var embed discordgo.MessageEmbed
	err := nil
	subs := []string{"dankmemes", "funny", "memes", "dank_meme", "comedyheaven", "CyanideandHappiness", "therewasanattempt", "wholesomememes", "instant_regret"}
	imageEndings := [6]string{".jpg", ".png", ".jpeg", ".gif", ".gifv", "gfycat"}
	limit := 100
	score, url, title, nsfw := GetPost(subs, limit)
	if strings.ContainsAny(url, imageEndings) {
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
	} else {
		err = "Url endings not in acceptable range."
	}

	return embed, nsfw, err
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
