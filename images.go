package main

// https://github.com/vivithemage/mrisa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type imagePayload struct {
	ImageURL      string `json:"image_url"`
	ResizedImages bool   `json:"resized_images"`
}

type returnPayload struct {
	Links         []string `json:"links"`
	Descriptions  []string `json:"descriptions"`
	SimilarImages []string `json:"similar_images"`
	Titles        []string `json:"titles"`
}

func imageRedditSearch(url string) string {
	var payload imagePayload
	var parsedBody returnPayload
	payload.ImageURL = url
	payload.ResizedImages = false
	var returnURL string
	var oldest int64

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Post(mrisaAddress, "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
		return "nil"
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &parsedBody)

	if err != nil {
		fmt.Println(err)
	}
	urls := parsedBody.Links
	titles := parsedBody.Descriptions
	oldest = 0
	for i := 0; i < len(urls); i++ {
		if strings.Contains(urls[i], "www.reddit.com/r/") {
			title := titles[i]
			dash := strings.Index(title, "-")
			newTitle := title[dash+2:]
			dash = strings.Index(newTitle, "-")
			timeStamp := newTitle[:dash-2]
			testTime := timeStrToSeconds(timeStamp)
			if testTime >= oldest {
				oldest = testTime
				returnURL = urls[i]
			}
		}
	}

	return returnURL

}

func imageSearch(url string) []string {
	var payload imagePayload
	var parsedBody returnPayload
	payload.ImageURL = url
	payload.ResizedImages = false

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Post(mrisaAddress, "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &parsedBody)

	if err != nil {
		fmt.Println(err)
	}
	urls := parsedBody.Links

	return urls
}

func imageSearchCommand(discord *discordgo.Session, channel string) {
	searchURL := ""
	extensions := []string{".jpg", ".png", ".jpeg"}
	discord.ChannelMessageSend(channel, "Searching Reddit...")
	lastMessages, _ := discord.ChannelMessages(channel, 100, "", "", "")
	for _, message := range lastMessages {
		attachmentsLength := len(message.Attachments)
		currentIndex := attachmentsLength - 1
		if attachmentsLength > 0 {
			for currentIndex >= 0 {
				last := message.Attachments[currentIndex]
				if ContainsAnySubstring(last.ProxyURL, extensions) {
					searchURL = last.ProxyURL
					break
				} else {
					currentIndex--
				}
			}
		}

		if searchURL != "" {
			break
		}

		if ContainsAnySubstring(message.Content, extensions) && !message.Author.Bot {
			httpIndex := strings.Index(message.Content, "http")
			startLink := message.Content[httpIndex:len(message.Content)]
			spaceIndex := strings.Index(startLink, " ")
			if spaceIndex != -1 {
				searchURL = startLink[0:spaceIndex]
				break
			} else {
				searchURL = startLink
			}
		}
	}

	if searchURL == "" {
		discord.ChannelMessageSend(channel, "404: URL not found. If you think that this is a mistake, post on our Github issues page along with appropriate screenshots and information. https://github.com/chand1012/Discord-Quick-Meme/issues")
		return
	}

	url := imageRedditSearch(searchURL)
	if url == "" {
		discord.ChannelMessageSend(channel, "Couldn't find anything on Reddit, searching the web....")
		urls := imageSearch(searchURL)
		if urls == nil {
			discord.ChannelMessageSend(channel, "500: Error connecting to image search service. If this persists, report at the Github issues page found here: https://github.com/chand1012/Discord-Quick-Meme/issues")
			return
		}
		printstr := "Found " + strconv.Itoa(len(urls)) + " results:\n"
		for _, link := range urls {
			printstr += link + "\n"
		}
		discord.ChannelMessageSend(channel, printstr)
	} else if url == "nil" {
		discord.ChannelMessageSend(channel, "500: Error connecting to image search service. If this persists, report at the Github issues page found here: https://github.com/chand1012/Discord-Quick-Meme/issues")
	} else {
		discord.ChannelMessageSend(channel, "I think I found the meme: \n"+url)
	}
}
