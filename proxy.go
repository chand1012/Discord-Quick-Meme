package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

func getFileType(contentURL string) string {
	reversedURL := reverseString(contentURL)
	endingIndex := strings.Index(reversedURL, ".")
	reversedEnding := []rune(reversedURL)[0:endingIndex]
	return reverseString(fmt.Sprint(reversedEnding))
}

func supportedType(contentURL string) bool {
	fileType := getFileType(contentURL)
	supported := []string{"jpg", "jpeg", "png", "gif", "gifv", "svg"}
	return stringInSlice(fileType, supported)
}

// will be finished later, comment out for now
func proxyEmbedSendRoutine(discord *discordgo.Session, channel string, sub string, title string, contentURL string, score int32) {
	proxyBase := "https://image-proxy.chand1012.workers.dev/"
	imageURL := proxyBase + contentURL
	rand.Seed(time.Now().Unix())
	randColor := rand.Intn(0xffffff)
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       randColor,
		Description: "From r/" + sub + "\n Score: " + humanize.Comma(int64(score)),
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     title,
	}
	_, err := discord.ChannelMessageSendEmbed(channel, embed)
	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	}
}

// This only applies to images at the moment, so the post should be checked beforehand
// Also should be checked if they have this setting set.
func uploadEmbedSendRoutine(discord *discordgo.Session, channel string, sub string, title string, contentURL string, score int32) {
	req, err := http.NewRequest("GET", contentURL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	}
	defer resp.Body.Close()

	imageName := path.Base(req.URL.Path)

	rand.Seed(time.Now().Unix())
	randColor := rand.Intn(0xffffff)

	message := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + imageName,
			},
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       randColor,
			Description: "From r/" + sub + "\n Score: " + humanize.Comma(int64(score)),
			Timestamp:   time.Now().Format(time.RFC3339),
			Title:       title,
		},
		Files: []*discordgo.File{
			{
				Name:   imageName,
				Reader: resp.Body,
			},
		},
	}
	_, err = discord.ChannelMessageSendComplex(channel, message)

	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	}
}

func uploadSendRoutine(discord *discordgo.Session, channel string, sub string, title string, contentURL string, score int32) {
	req, err := http.NewRequest("GET", contentURL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	}
	defer resp.Body.Close()

	imageName := path.Base(req.URL.Path)

	message := &discordgo.MessageSend{
		Content: "From r/" + sub + "\n" + title + "\nScore: " + humanize.Comma(int64(score)),
		Files: []*discordgo.File{
			{
				Name:   imageName,
				Reader: resp.Body,
			},
		},
	}

	_, err = discord.ChannelMessageSendComplex(channel, message)

	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	}
}

func proxySendRoutine(discord *discordgo.Session, channel string, sub string, title string, contentURL string, score int32) {
	proxyBase := "https://image-proxy.chand1012.workers.dev/"
	imageURL := proxyBase + contentURL
	_, err := discord.ChannelMessageSend(channel, "From r/"+sub+"\n"+title+"\n"+imageURL+"\nScore: "+humanize.Comma(int64(score)))
	if err != nil {
		fmt.Println("Error posting to channel:", err.Error())
	}
}
