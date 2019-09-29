package main

import (
	"strings"

	"github.com/go-redis/redis"
)

// this will handle all redis related commands

// GetBannedSubreddits gets a list of banned subs from redis
func GetBannedSubreddits(channel string) ([]string, error) {
	address, password, db, err := redisExtract("data.json")
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()
	rawValues, err := redisClient.Get(channel).Result()
	if err != nil {
		return nil, err
	}

	values := strings.Split(rawValues, " ")

	return values, err
}

//AppendBannedSubreddit appends a banned subreddit to the list for that channel
func AppendBannedSubreddit(channel string, subreddit string) error {
	var isContained bool
	address, password, db, err := redisExtract("data.json")
	if err != nil {
		return err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()

	values, err := redisClient.Get(channel).Result()
	if err != nil {
		return err
	}

	if values == "" {
		values = subreddit
		isContained = false
	} else if strings.Contains(values, subreddit) {
		isContained = true
	} else {
		values = values + " " + subreddit
		isContained = false
	}

	if !isContained {
		err = redisClient.Set(channel, values, 0).Err()
	} else {
		err = nil
	}

	return err

}

func UnbanSubreddit(channel string, subreddit string) error {
	var isContained bool
	address, password, db, err := redisExtract("data.json")
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()

	values, err := redisClient.Get(channel).Result()
	if err != nil {
		return err
	}

	if strings.Contains(values, " ") {
		if strings.Contains(values, subreddit) {
			values = strings.Replace(values, " "+subreddit, "", -1)
			isContained = true
		} else {
			isContained = false
		}
	} else {
		if strings.Contains(values, subreddit) {
			values = ""
			isContained = true
		} else {
			isContained = false
		}
	}

	if isContained {
		err = redisClient.Set(channel, values, 0).Err()
	} else {
		err = nil
	}

	return err

}
