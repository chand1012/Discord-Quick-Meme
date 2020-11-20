package main

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// var (
// 	ServerMap map[string]string
// 	NSFWMap   map[string]bool
// )

func TestChannelDataGetter(t *testing.T) {

	// these for testing the function
	ServerMap = make(map[string]string)
	NSFWMap = make(map[string]bool)

	testChannel := os.Getenv("TEST_CHANNEL")
	testGuild := os.Getenv("TEST_GUILD")
	token := os.Getenv("DISCORD_TOKEN")

	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		t.Errorf("There was an error creating the discord bot: %v", err)
	}

	defer discord.Close()

	// time the test because GetChannelData should be ran twice. The second time should be faster
	firstStart := GetMillis()

	name, nsfw := GetChannelData(discord, testChannel, testGuild)

	firstEnd := GetMillis()
	firstTime := firstEnd - firstStart

	if name != "bottest" || nsfw != false { // I wrote the comparison like this to be extra readable
		t.Errorf("There was an error getting channel name and NSFW status from DB. Expected 'bottest' and 'false', got '%s' and '%t'.", name, nsfw)
	}

	secondStart := GetMillis()

	name, nsfw = GetChannelData(discord, testChannel, testGuild)

	secondEnd := GetMillis()
	secondTime := secondEnd - secondStart

	if name != "bottest" || nsfw != false { // I wrote the comparison like this to be extra readable
		t.Errorf("There was an error getting channel name and NSFW status from cache. Expected 'bottest' and 'false', got '%s' and '%t'.", name, nsfw)
	}

	if secondTime >= firstTime {
		t.Errorf("Time to get from cache longer than or equal to database.")
	}

	// cleanup
	err = RemoveChannelFromDB(testChannel)

	if err != nil {
		t.Errorf("Error removing channel from database: %v", err)
	}
}
