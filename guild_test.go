package main

import (
	"math/rand"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestServerSettings(t *testing.T) {

	SettingsMap = make(map[string]guildSettings)

	testGuild := os.Getenv("TEST_GUILD")
	token := os.Getenv("DISCORD_TOKEN")

	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		t.Errorf("There was an error creating the discord bot: %v", err)
	}

	defer discord.Close()

	testOneStart := GetMillis()
	_, err = getServerSettings(discord, testGuild)

	testOneEnd := GetMillis()
	testOne := testOneEnd - testOneStart

	if err != nil {
		t.Errorf("There was an error getting server settings for guild %s: %v", testGuild, err)
	}

	testTwoStart := GetMillis()

	_, err = getServerSettings(discord, testGuild)

	testTwoEnd := GetMillis()
	testTwo := testTwoEnd - testTwoStart

	if testOne <= testTwo {
		t.Errorf("Cache get took longer than the database get. Cache time: %d; Database time: %d", testTwo, testOne)
	}

	proxyEnable := (rand.Intn(2) == 1)
	proxyMode := int8(rand.Intn(100))

	updateSettingsCache(testGuild, proxyEnable, proxyMode)

	settings, err := getServerSettings(discord, testGuild)

	if err != nil {
		t.Errorf("There was an error getting server settings for guild %s: %v", testGuild, err)
	}

	if settings.Proxy != proxyEnable {
		t.Errorf("Proxy enable settings do not match. Expected %t, got %t.", proxyEnable, settings.Proxy)
	}

	if settings.ProxyMode != proxyMode {
		t.Errorf("Proxy mode settings do not match. Expected %d, got %d.", proxyMode, settings.ProxyMode)
	}

}
