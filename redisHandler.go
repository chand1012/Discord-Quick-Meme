package main

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis"
)

// this will handle all redis related commands

func redisSave() error {
	address, password, db, err := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()

	cmd := redis.NewStringCmd("save")
	err = redisClient.Process(cmd)
	return err
}

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
	if err != nil {
		fmt.Println(err)
	}
	values := strings.Split(rawValues, " ")

	return values, err
}

// SubIsBanned checks if a subreddit is banned on that channel
func SubIsBanned(channel string, subreddit string) (bool, error) {
	bannedSubs, err := GetBannedSubreddits(channel)
	for _, sub := range bannedSubs {
		if subreddit == sub {
			return true, err
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return false, err
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
	if err != nil {
		fmt.Println(err)
	}
	if values == "" || strings.Replace(values, " ", "", -1) == "" {
		values = subreddit
		err = redisClient.Set(channel, values, 0).Err()
	} else if !strings.Contains(values, subreddit) {
		values = values + " " + subreddit
		err = redisClient.Set(channel, values, 0).Err()
	} else {
		err = nil
	}
	if err != nil {
		fmt.Println(err)
	}
	err = redisSave()
	return err

}

// UnbanSubreddit removes a subreddit from the redis banned servers
func UnbanSubreddit(channel string, subreddit string) error {
	var isContained bool
	address, password, db, err := redisExtract("data.json")
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
		if err != nil {
			fmt.Println(err)
		}
		err = redisSave()
	} else {
		err = nil
	}
	if err != nil {
		fmt.Println(err)
	}

	return err

}

func saveCommonSubsRedis() error {
	var redisSubs []string
	for sub, count := range CommonSubs {
		if count >= 10 {
			redisSubs = append(redisSubs, sub)
		}
	}
	CommonSubsCounter = uint8(len(redisSubs))
	if redisSubs == nil {
		return nil
	}
	list := strings.Join(redisSubs, " ")
	address, password, db, err := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	defer redisClient.Close()
	err = redisClient.Set("commonSubs", list, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	go redisSave()
	return nil
}

func getCommonSubsRedis() ([]string, error) {
	var redisSubs []string
	list, err := getCommonSubsRedisRaw()
	if list == "" {
		return nil, err
	}
	redisSubs = strings.Split(list, " ")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return redisSubs, err
}

func getCommonSubsRedisRaw() (string, error) {
	address, password, db, err := redisExtract("data.json")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer redisClient.Close()
	list, err := redisClient.Get("commonSubs").Result()
	if err != nil {
		fmt.Println(err)
	}
	return list, err
}
