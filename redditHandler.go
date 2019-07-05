package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
)

// QuickPost variant of reddit.Post but designed to be pulled from ram on the fly
type QuickPost struct {
	Title     string
	Score     int32
	Content   string
	Nsfw      bool
	Permalink string
}

// GetFromCache pulls post from the cache
func GetFromCache(sub string) ([]QuickPost, bool) {
	var success bool
	var returnPosts []QuickPost
	returnPosts, success = PostCache[sub]
	return returnPosts, success
}

// AddToCache adds post to cache
func AddToCache(sub string, addPosts []QuickPost) {
	PostCache[sub] = addPosts
}

// ClearCache clears the cache
func ClearCache() {
	for k := range PostCache {
		delete(PostCache, k)
	}
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

	var urlItems []string
	var sub string

	var success bool

	var s int

	sub = subs[rand.Intn(len(subs))]
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
			switch mode {
			case "link":
				gotPost = QuickPost{
					Title:     post.Title,
					Score:     post.Score,
					Content:   post.URL,
					Nsfw:      post.NSFW,
					Permalink: post.Permalink,
				}
			case "text":
				gotPost = QuickPost{
					Title:     post.Title,
					Score:     post.Score,
					Content:   post.SelfText,
					Nsfw:      post.NSFW,
					Permalink: post.Permalink,
				}
			case "media":
				urlItems = []string{".jpg", ".png", ".jpeg", "gfycat", "youtube", "youtu.be", "gif", "gifv"}
				if !strings.Contains(post.URL, "v.redd.it") && ContainsAnySubstring(post.URL, urlItems) {
					gotPost = QuickPost{
						Title:     post.Title,
						Score:     post.Score,
						Content:   post.URL,
						Nsfw:      post.NSFW,
						Permalink: post.Permalink,
					}
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
		s := rand.Intn(len(cachePosts))
		returnPost = cachePosts[s]
	}
	return returnPost, sub
}
