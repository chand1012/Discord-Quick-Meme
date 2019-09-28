package main

import (
	"strings"

	"github.com/go-redis/redis"
)

// this will handle all redis related commands

// GetBannedSubreddits gets a list of banned subs from redis
func GetBannedSubreddits(channel string) ([]string, error) {
	address, password, db, err := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()
	rawValues, err := redisClient.Get(channel).Result()

	values := strings.Split(rawValues, " ")

	return values, err
}

//AppendBannedSubreddit appends a banned subreddit to the list for that channel
func AppendBannedSubreddit(channel string, subreddit string) error {
	address, password, db, err := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()

	values, err := redisClient.Get(channel).Result()

	if values == "" {
		values = subreddit
	} else {
		values = values + " " + subreddit
	}
	err = redisClient.Set(channel, values, 0).Err()

	return err

}
