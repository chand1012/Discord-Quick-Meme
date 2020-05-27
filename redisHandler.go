package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
)

//initializes the redis client
func initRedis() *redis.Client {
	address, password, db := getRedisEnv()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return redisClient
}

// this is because we are going to use a separate database within redis for all queue data
func initRedisOverride(oAddress string, oPassword string, oDB int) *redis.Client {
	address, password, db := getRedisEnv()

	if oAddress != "" {
		address = oAddress
	}
	if oPassword != "" {
		password = oPassword
	}
	if oDB != -1 {
		db = oDB
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return redisClient
}

// QueueObj data structure for the queue
type QueueObj struct {
	Interval   string   `json:"interval"`
	Time       int64    `json:"time"`
	Type       string   `json:"type"`
	SubReddits []string `json:"subreddit"`
	NSFW       bool     `json:"nsfw"`
}

func setRedisQueue(channel string, timeframe string, postType string, subs []string, nsfw bool) error {
	var redisQueue QueueObj
	var redisDB int

	redisDB = 1
	if RunMode == "dev" {
		redisDB = 2
	}

	redisClient := initRedisOverride("", "", redisDB)
	defer redisClient.Close()

	redisQueue.Interval = timeframe
	redisQueue.SubReddits = subs
	redisQueue.Type = postType
	redisQueue.NSFW = nsfw
	redisQueue.Time = 0

	jsonString, err := json.Marshal(redisQueue)

	if err != nil {
		return err
	}

	err = redisClient.Set(channel, string(jsonString), 0).Err()
	return err
}

func setRedisQueueRaw(redisQueue QueueObj, channel string) error {
	var redisDB int

	redisDB = 1
	if RunMode == "dev" {
		redisDB = 2
	}

	redisClient := initRedisOverride("", "", redisDB)
	defer redisClient.Close()

	jsonString, err := json.Marshal(redisQueue)
	if err != nil {
		return err
	}

	QueueState[channel] = false

	err = redisClient.Set(channel, string(jsonString), 0).Err()
	return err
}

func getRedisQueue(channel string) (QueueObj, error) {
	var redisQueue QueueObj
	var redisDB int

	redisDB = 1
	if RunMode == "dev" {
		redisDB = 2
	}

	redisClient := initRedisOverride("", "", redisDB)
	defer redisClient.Close()

	value, err := redisClient.Get(channel).Result()

	if err != nil {
		return QueueObj{}, err
	}

	json.Unmarshal([]byte(value), &redisQueue)

	return redisQueue, err
}

func getAllQueueChannels() ([]string, error) {
	var redisDB int

	redisDB = 1
	if RunMode == "dev" {
		redisDB = 2
	}

	redisClient := initRedisOverride("", "", redisDB)
	defer redisClient.Close()
	keys, _ := redisClient.Do("keys", "*").Result()

	return interfaceToStringSlice(keys), nil
}

func redisDelete(key string, database int) error {
	redisClient := initRedisOverride("", "", database)
	defer redisClient.Close()

	err := redisClient.Get(key).Err()

	if err != nil {
		return err
	}

	return redisClient.Del(key).Err()
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
	if err == redis.Nil {
		return nil, nil
	}
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
func AppendBannedSubreddit(channel string, subreddit string) error {

	redisClient := initRedis()
	defer redisClient.Close()

	values, err := redisClient.Get(channel).Result()
	if err == redis.Nil {
		values = ""
	} else if err != nil {
		fmt.Println("Error getting redis values: ", err)
		return err
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
		fmt.Println("Error setting redis values: ", err)
		return err
	}
	err = redisSave()
	if err != nil {
		fmt.Println("Error saving redis: ", err)
		return nil
	}

	return nil

}

// UnbanSubreddit removes a subreddit from the redis banned servers
func UnbanSubreddit(channel string, subreddit string) error {
	var isContained bool

	redisClient := initRedis()

	defer redisClient.Close()

	values, err := redisClient.Get(channel).Result()
	if err == redis.Nil {
		values = ""
	} else if err != nil {
		fmt.Println("Error getting redis values: ", err)
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
			fmt.Println("Error setting redis values: ", err)
			return err
		}
		err = redisSave()
		if err != nil {
			fmt.Println("Error saving redis: ", err)
			return nil
		}
	}

	return nil

}
