package main

import (
	"fmt"

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

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
