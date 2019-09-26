package main

import (
	"strings"

	"github.com/go-redis/redis"
)

// this will handle all redis related commands

func GetBannedSubreddits(channel string) []string {
	address, password, db := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()
	raw_values, err := redisClient.Get(channel).Result()

	values := strings.Split(raw_values, ",")

	return values
}

func AppendBannedSubreddit(channel string, subreddit string) error {
	address, password, db := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()

	raw_values, err := redisClient.Get(channel).Result()

	var values []string
	var appendString string
	var resultString string

	if strings.Contains(raw_values, ",") {
		values = strings.Split(raw_values, ",")
		appendString = ","
	} else {
		values = []string{raw_values}
		appendString = ""
	}

	for _, sub := range values {
		resultString += (appendString + sub)
	}

}
