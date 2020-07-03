package main

import (
	"fmt"
	"net/url"
	"strings"
)

// loop routine that gets posts from the cache.
func getPostLoop(subs []string, postType string, channel string, channelNsfw bool, sort string) (string, int32, string, string, bool, bool) {
	var returnPost QuickPost
	var sub string
	var score int32
	var contentURL string
	var title string
	var nsfw bool

	bannedToggle := false
	toggled := false
	bannedSubs, err := GetAllBannedSubs(channel)
	if err != nil {
		fmt.Println("Error getting banned subreddits: ", err)
	}
	for i := 0; i < 10; i++ {
		returnPost, sub = GetPost(subs, 100, sort, postType)
		blacklisted := CheckBlacklist(channel, returnPost)
		banned := stringInSlice(sub, bannedSubs)

		if !banned {
			LastPost[channel] = returnPost
		}

		score = returnPost.Score
		contentURL = returnPost.Content
		title = returnPost.Title
		nsfw = returnPost.Nsfw

		if nsfw {
			fmt.Println("Post is NSFW.")
		}

		if strings.Contains(contentURL, "imgur") {
			if !strings.Contains(contentURL, "i.imgur") && !strings.Contains(contentURL, "/a/") && !strings.Contains(contentURL, "gif") {
				contentURL = getImgurDirectLink(contentURL)
			}
		}

		// complicated but it makes sense ish
		if (channelNsfw && !blacklisted && !banned) || (channelNsfw && !nsfw && !blacklisted && !banned) {
			toggled = true
			AddToBlacklist(channel, returnPost)
			break
		} else if !channelNsfw && !nsfw && !blacklisted && !banned {
			toggled = true
			break
		} else {
			if !blacklisted && !banned {
				fmt.Println("Channel is not NSFW but post is NSFW, retrying....")
			} else if banned {
				fmt.Println("Channel banned sub " + sub + ", retrying...")
				if i == 9 {
					bannedToggle = true
				}
			}
		}
	}
	return sub, score, contentURL, title, toggled, bannedToggle

}

func getImgurDirectLink(contentURL string) string {
	parsedURL, err := url.Parse(contentURL)
	if err != nil {
		fmt.Println(err)
		return contentURL
	}
	imgurDirect := "https://i.imgur.com"
	// use jpeg because its small and works
	return imgurDirect + parsedURL.Path + ".jpg"
}
