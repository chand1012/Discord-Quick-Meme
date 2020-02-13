package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
)

// GuessPostType get the type of the post so a
func GuessPostType(post *reddit.Post) string {
	selfText := post.SelfText
	urlContent := post.URL
	if len(selfText) == 0 {
		urlItems := []string{".jpg", ".png", ".jpeg", "gfycat", "youtube", "youtu.be", "gif", "gifv"}
		if !strings.Contains(urlContent, "v.redd.it") && ContainsAnySubstring(urlContent, urlItems) {
			return "media"
		}
		return "link"

	}
	return "text"
}

// PingReddit tests reddit connection
func PingReddit() error {
	bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
	_, err = bot.Listing("/r/all", "")
	return err
}

// GetPost gets reddit posts
func GetPost(subs []string, limit int, sort string, mode string) (QuickPost, string) {
	var gottenPosts []QuickPost
	var cachePosts []QuickPost
	var gotPost QuickPost
	var returnPost QuickPost

	var subList []string
	var sub string

	var success bool

	var s int

	subList = getAllSubsFromMap()
	sub = subs[rand.Intn(len(subs))]
	go updateCommonSubCounter(sub)
	cachePosts, success = GetFromCache(sub)
	now := time.Now().Unix()
	if now >= CacheTime {
		fmt.Println("Clearing Cache...")
		ClearCache()
		success = false
		CacheTime = time.Now().Unix() + 3600
		fmt.Println("New cache time is " + strconv.FormatInt(CacheTime, 10))
	}
	switch {
	case !success:
		starttime := GetMillis()
		fmt.Println("Adding r/" + sub + " to cache.")
		bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
		if err != nil {
			panic(err)
		}
		rand.Seed(time.Now().Unix())
		harvest, err := bot.Listing("/r/"+sub+"/"+sort, "")
		lengthPosts := len(harvest.Posts)
		if lengthPosts < limit {
			limit = lengthPosts
		}
		for _, post := range harvest.Posts[:limit] {
			mode := GuessPostType(post)
			switch {
			case mode == "link" || mode == "media":
				gotPost = QuickPost{
					Title:     post.Title,
					Score:     post.Score,
					Content:   post.URL,
					Nsfw:      post.NSFW,
					Permalink: post.Permalink,
					Sub:       getSubFromPermalink(post.Permalink),
				}
			case mode == "text":
				gotPost = QuickPost{
					Title:     post.Title,
					Score:     post.Score,
					Content:   post.SelfText,
					Nsfw:      post.NSFW,
					Permalink: post.Permalink,
					Sub:       getSubFromPermalink(post.Permalink),
				}
			}

			gottenPosts = append(gottenPosts, gotPost)
		}
		gottenLength := len(gottenPosts)
		if gottenLength == 0 {
			returnPost = QuickPost{
				Title:     "ERROR: Sub seems to be empty or does not exist.",
				Score:     0,
				Content:   "",
				Nsfw:      false,
				Permalink: "/r/" + sub + "/",
			}
			fmt.Println("Nothing to cache! Discarding....")
		} else if ContainsAnySubstring(sub, subList) {
			s = rand.Intn(gottenLength)
			returnPost = gottenPosts[s]
			if !CachePopulating {
				go PopulateCache()
			}
		} else {
			AddToCache(sub, gottenPosts)
			CacheTime = time.Now().Unix() + 3600
			s = rand.Intn(gottenLength)
			returnPost = gottenPosts[s]
			endtime := GetMillis()
			t := endtime - starttime
			fmt.Println("Took " + strconv.FormatInt(t, 10) + "ms to add to cache.")
		}
	case success:
		fmt.Println("Found r/" + sub + " in cache.")
		minScore := MinScore(cachePosts)
		for i := 0; i < len(cachePosts); i++ {
			s := rand.Intn(len(cachePosts))
			returnPost = cachePosts[s]
			if returnPost.Score >= minScore {
				break
			}
		}
	}
	return returnPost, getSubFromPermalink(returnPost.Permalink)
}

// MinScore Formula for calculating effective minimum score.
func MinScore(posts []QuickPost) int32 {
	var total int32
	var n int
	n = len(posts)
	for _, post := range posts {
		total += post.Score
	}
	avg := total / int32(n)
	return avg / 2
}

// getSubFromPermalink gets a sub from the link to the post
func getSubFromPermalink(permalink string) string {
	var sub string
	linkArray := strings.Split(permalink, "/")
	sub = linkArray[2]
	return sub
}

// getSubsFromMap gets the subreddits from RAM instead of from disk
func getAllSubsFromMap() []string {
	var allSubs []string
	for _, value := range SubMap {
		for _, sub := range value {
			allSubs = append(allSubs, sub)
		}
	}
	return allSubs
}
