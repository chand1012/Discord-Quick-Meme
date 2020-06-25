package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
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

func proxySendRoutine(discord *discordgo.Session, channel string, sub string, title string, url string, score int32) {

}
