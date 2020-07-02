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

func getFileType(url string) string {
	reversedURL := reverseString(url)
	endingIndex := strings.Index(reversedURL, ".")
	reversedEnding := []rune(reversedURL)[0:endingIndex]
	return reverseString(string(reversedEnding))
}

func supportedType(url string) bool {
	fileType := getFileType(url)
	supported := []string{"jpg", "jpeg", "png", "gif", "gifv", "svg"}
	return stringInSlice(fileType, supported)
}

// will be finished later, comment out for now
func proxySendRoutine(discord *discordgo.Session, channel string, sub string, title string, url string, score int32) {
	proxyBase := "https://image-proxy.chand1012.workers.dev/"
	imageURL := proxyBase + url
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
func fileUploadRoutine(discord *discordgo.Session, channel string, sub string, title string, url string, score int32) {
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	}
	defer resp.Body.Close()

	imageName := path.Base(req.URL.Path)

	if err != nil {
		fmt.Println(err)
		errSendRoutine(discord, channel, err)
	}

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
			&discordgo.File{
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
