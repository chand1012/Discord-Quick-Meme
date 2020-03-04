package main

import (
	"strings"

	"github.com/go-redis/redis"
)

//initializes the redis client
func initRedis() *redis.Client {
	address, password, db, err := redisExtract("data.json")
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return redisClient
}

// redisSave saves redis cache to disk
func redisSave() error {
	redisClient := initRedis()
	defer redisClient.Close()

	cmd := redis.NewStringCmd("save")
	err := redisClient.Process(cmd)
	return err
}

// GetBannedSubreddits gets a list of banned subs from redis
func GetBannedSubreddits(channel string) ([]string, error) {
	redisClient := initRedis()
	defer redisClient.Close()
	rawValues, err := redisClient.Get(channel).Result()
	if err != nil {
		return nil, err
	}
	values := strings.Split(rawValues, " ")
	return values, err
}

// SubIsBanned checks if a subreddit is banned on that channel
// May be used later but not used at the moment, so commented out
// func SubIsBanned(channel string, subreddit string) (bool, error) {
// 	bannedSubs, err := GetBannedSubreddits(channel)
// 	for _, sub := range bannedSubs {
// 		if subreddit == sub {
// 			return true, err
// 		}
// 	}
// 	return false, err
// }

//AppendBannedSubreddit appends a banned subreddit to the list for that channel
func AppendBannedSubreddit(channel string, subreddit string) {

	redisClient := initRedis()
	defer redisClient.Close()

	values, err := redisClient.Get(channel).Result()
	errCheck("Error getting redis values", err, false)
	if values == "" || strings.Replace(values, " ", "", -1) == "" {
		values = subreddit
		err = redisClient.Set(channel, values, 0).Err()
	} else if !strings.Contains(values, subreddit) {
		values = values + " " + subreddit
		err = redisClient.Set(channel, values, 0).Err()
	} else {
		err = nil
	}
	errCheck("Error setting value '"+subreddit+"'", err, false)
	err = redisSave()
	errCheck("Error saving redis", err, false)

}

// UnbanSubreddit removes a subreddit from the redis banned servers
func UnbanSubreddit(channel string, subreddit string) {
	var isContained bool

	redisClient := initRedis()

	defer redisClient.Close()

	values, err := redisClient.Get(channel).Result()
	errCheck("Error getting redis values", err, false)

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
		errCheck("Error setting value '"+values+"'", err, false)
		err = redisSave()
		errCheck("Error saving redis", err, false)
	}

}

// saveCommonSubsRedis saves a subreddit to the common sub cache in redis
func saveCommonSubsRedis() error {
	var redisSubs []string
	for sub, count := range CommonSubs {
		if count >= 5 {
			redisSubs = append(redisSubs, sub)
		}
	}
	CommonSubsCounter = uint8(len(redisSubs))
	if redisSubs == nil {
		return nil
	}
	list := strings.Join(redisSubs, " ")

	redisClient := initRedis()
	defer redisClient.Close()

	err := redisClient.Set("commonSubs", list, 0).Err()
	if err != nil {
		return err
	}
	err = redisSave()
	return err
}

// getCommonSubsRedis gets all the commonly used subs as a list of strings
func getCommonSubsRedis() ([]string, error) {
	var redisSubs []string
	list, err := getCommonSubsRedisRaw()
	if list == "" || err != nil {
		return nil, err
	}
	redisSubs = strings.Split(list, " ")
	return redisSubs, err
}

// getCommonSubsRedisRaw gets commonly used subs as a space seperated string
func getCommonSubsRedisRaw() (string, error) {
	redisClient := initRedis()
	defer redisClient.Close()
	list, err := redisClient.Get("commonSubs").Result()
	return list, err
}
