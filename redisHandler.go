package main

import (
	"strings"

	"github.com/go-redis/redis"
)

// this will handle all redis related commands

// GetBannedSubreddits gets a list of banned subs from redis
func GetBannedSubreddits(channel string) []string {
	address, password, db := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()
	rawValues, err := redisClient.Get(channel).Result()

	values := strings.Split(rawValues, ",")

	return values
}

//AppendBannedSubreddit appends a banned subreddit to the list for that channel
func AppendBannedSubreddit(channel string, subreddit string) error {
	address, password, db := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()

	rawValues, err := redisClient.Get(channel).Result()

	var values []string
	var appendString string
	var resultString string

	if strings.Contains(rawValues, ",") {
		values = strings.Split(rawValues, ",")
		values = append(values, subreddit)
		appendString = ","
	} else {
		values = []string{rawValues, subreddit}
		appendString = ""
	}

	for _, sub := range values {
		resultString += (appendString + sub)
	}

	err = redisClient.Set(channel, resultString, 0).Err()

	return err

}
