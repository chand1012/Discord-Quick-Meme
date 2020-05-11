package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type agentFile struct {
	UserAgent    string `yaml:"user_agent"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

func generateAgentFile() error {
	var agent agentFile
	agent.UserAgent = "DiscordQuickMemeBot"
	agent.ClientID = os.Getenv("REDDIT_ID")
	agent.ClientSecret = os.Getenv("REDDIT_SECRET")

	if agent.ClientID == "" || agent.ClientSecret == "" {
		fmt.Println("Cannot continue, Reddit client ID or Secret not set.")
		os.Exit(1)
	}

	rawYAML, err := yaml.Marshal(&agent)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./agent.yml", rawYAML, 0644)

	return err
}

func getRedisEnv() (string, string, int) {
	var redisDB int
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPasswd := os.Getenv("REDIS_PASSWORD")
	redisDBRaw := os.Getenv("REDIS_DB")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	if redisDBRaw == "" {
		redisDB = 0
	} else {
		redisDB, _ = strconv.Atoi(redisDBRaw)
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
