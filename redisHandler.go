package main

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis"
)

//initializes the redis client
func initRedis() *redis.Client {
	address, password, db, err := redisExtract("data.json")
	if err != nil {
		panic(err) // cannot launch without this
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return redisClient
}

// this is because we are going to use a seperate database within redis for all queue data
func initRedisOverride(oAddress string, oPassword string, oDB int) *redis.Client {
	address, password, db, err := redisExtract("data.json")
	if err != nil {
		panic(err)
	}
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
	if err == redis.Nil {
		list = ""
		err = nil
	}
	return list, err
}
