package main

import (
	"lib/jsonHandler"
	"lib/redditHandler"
	"fmt"
	"github.com/turnage/graw"
)
var (
	commandPrefix string
	botId		  string
)

func main() {
	keys, err := jsonHandler.extract("data.json")
	errCheck("Error opening key file", err)
	discord, err := discordgo.New("Bot " + keys.botId)
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
	err = discord.updateStatus(0, "A Reddit meme bot.")
	if err != nil {
		fmt.Println("Error attempting to set the status.")
	}
	servers := discord.State.Guilds
	fmt.Printf("Discord-Quick-Meme has started on %d servers", len(servers))
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		return
	}
	// this will be epanded upon in the near future
	content := message.Content
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
