package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type agentFile struct {
	UserAgent    string `yaml:"user_agent"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
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

func getRedisEnv() (string, string, int) {
	var redisDB int
	var err error
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPasswd := os.Getenv("REDIS_PASSWORD")
	redisDBRaw := os.Getenv("REDIS_DB")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	if redisDBRaw == "" {
		redisDB = 0
	} else {
		redisDB, err = strconv.Atoi(redisDBRaw)
		if err != nil {
			redisDB = 0
		}
	}

	return redisAddr, redisPasswd, redisDB
}

func getMRISAEnv() string {
	mrisa := os.Getenv("MRISA")
	if mrisa == "" {
		mrisa = "http://192.168.1.2:5000/search"
	}
	return mrisa
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
