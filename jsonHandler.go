package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// this is for subreddit extraction. There will be one attribute for "memes" and the like
type subJSON struct {
	Memes      []string `json:"memes"`
	Text       []string `json:"text"`
	Hentai     []string `json:"hentai"`
	News       []string `json:"news"`
	FiftyFifty []string `json:"fiftyfifty"`
}

func subExtract(filename string) map[string][]string {
	jsonfile, err := os.Open(filename)
	if err != nil {
		panic(err) // these are situations that the bot cannot run if it errors out
	}

	defer jsonfile.Close()

	rawjson, err := ioutil.ReadAll(jsonfile)

	if err != nil {
		panic(err)
	}

	var subJSON subJSON

	json.Unmarshal(rawjson, &subJSON)

	subMap := make(map[string][]string)

	subMap["memes"] = subJSON.Memes
	subMap["text"] = subJSON.Text
	subMap["hentai"] = subJSON.Hentai
	subMap["news"] = subJSON.News
	subMap["fiftyfifty"] = subJSON.FiftyFifty

	return subMap
}

func getAllSubs(filename string) []string {
	var subs []string
	subMap := subExtract(filename)
	for key := range subMap {
		subsFromMap := subMap[key]
		for i := 0; i < len(subsFromMap); i++ {
			subs = append(subs, subsFromMap[i])
		}
	}
	return subs
}
