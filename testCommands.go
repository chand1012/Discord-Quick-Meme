package main

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// default case for the quickmeme command
func quickMemeDefault(discord *discordgo.Session, channel string) {
	servers := discord.State.Guilds
	//userCount := getNumberOfUsers(discord)
	msg := "Discord-Quick-Meme is active and ready on " + strconv.Itoa(len(servers)) + " servers."
	fmt.Println(msg)
	discord.ChannelMessageSend(channel, msg)
}
