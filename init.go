package main

import (
	"fmt"
	"os"
	"strings"
)

type agentFile struct {
	UserAgent    string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

func getRedditEnv() agentFile {
	var agent agentFile
	agent.UserAgent = "DiscordQuickMemeBot"
	agent.ClientID = os.Getenv("REDDIT_ID")
	agent.ClientSecret = os.Getenv("REDDIT_SECRET")

	if agent.ClientID == "" || agent.ClientSecret == "" {
		fmt.Println("Cannot continue, Reddit client ID or Secret not set.")
		os.Exit(1)
	}

	return agent
}

func getDataEnv() (string, string, string, []string) { // discord token, topgg key, mode, comma seperated list of admin ids
	token := os.Getenv("DISCORD_TOKEN")
	topKey := os.Getenv("TOPGG")
	mode := os.Getenv("MODE")
	adminsRaw := os.Getenv("ADMINS")

	admins := strings.Split(adminsRaw, ",")

	if token == "" {
		fmt.Println("Cannot continue, no Discord token specified.")
		os.Exit(1)
	}
	if mode == "" {
		mode = "prod"
	}
	return token, topKey, mode, admins
}
