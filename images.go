package main

// https://github.com/vivithemage/mrisa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// payload structs for sending and receiving search info
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

// MRISA search that only searches reddit
func imageRedditSearch(contentURL string) string {
	var returnURL string
	var oldest int64

	urls, titles := imageSearch(contentURL)

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

// the MRISA search function
func imageSearch(contentURL string) ([]string, []string) {
	var payload imagePayload
	var parsedBody returnPayload
	payload.ImageURL = contentURL
	payload.ResizedImages = false

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Error processing JSON data for payload: ", err)
		return nil, nil
	}

	mrisaAddress := os.Getenv("MRISA")
	if mrisaAddress == "" {
		mrisaAddress = "http://127.0.0.1:5000/search"
	}

	resp, err := http.Post(mrisaAddress, "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("Error connecting to MRISA server", err)
		return nil, nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading JSON data from response: ", err)
		return nil, nil
	}

	json.Unmarshal(body, &parsedBody)

	urls := parsedBody.Links
	titles := parsedBody.Descriptions
	return urls, titles
}

// The command that executes the above function and send the info to the channel
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

	contentURL := imageRedditSearch(searchURL)
	if contentURL == "" {
		discord.ChannelMessageSend(channel, "Couldn't find anything on Reddit, searching the web....")
		urls, _ := imageSearch(searchURL)
		if urls == nil {
			discord.ChannelMessageSend(channel, "500: Error connecting to image search service. If this persists, report at the Github issues page found here: https://github.com/chand1012/Discord-Quick-Meme/issues")
			return
		}
		printstr := "Found " + strconv.Itoa(len(urls)) + " results:\n"
		for _, link := range urls {
			printstr += link + "\n"
		}
		discord.ChannelMessageSend(channel, printstr)
	} else if contentURL == "nil" {
		discord.ChannelMessageSend(channel, "500: Error connecting to image search service. If this persists, report at the Github issues page found here: https://github.com/chand1012/Discord-Quick-Meme/issues")
	} else {
		discord.ChannelMessageSend(channel, "I think I found the meme: \n"+contentURL)
	}
}

// func getImgurDirectLink(contentURL string) string {

// }
