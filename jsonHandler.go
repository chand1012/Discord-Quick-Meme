package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Keys type for json keys
type Keys struct {
	BotID  string   `json:"token"`
	Admins []string `json:"admin"`
}

func jsonExtract(filename string) (string, []string, error) {
	jsonfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonfile.Close()

	rawjson, _ := ioutil.ReadAll(jsonfile)

	var keys Keys

	json.Unmarshal(rawjson, &keys)

	return keys.BotID, keys.Admins, err
}

// BanEntry an entry for a ban
type BanEntry struct {
	Subreddit string   `json:"subreddit"`
	Channels  []string `json:"channels"`
}

// GetBannedSubreddits gets all of the subreddits that are banned on the list of channels
func GetBannedSubreddits(filename string) map[string][]string, error {
	returnMap := make(map[string][]string)
	jsonfile, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer jsonfile.Close()

	rawjson, _ := ioutil.ReadAll(jsonfile)

	var entries []BanEntry

	json.Unmarshal(rawjson, &entries)

	for _, entry := range entries {
		returnMap[entry.Subreddit] = entry.Channels
	}

	return returnMap, err

}
