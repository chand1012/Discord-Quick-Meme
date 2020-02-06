package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	Sub       string
}

// GetFromCache pulls post from the cache
func GetFromCache(sub string) ([]QuickPost, bool) {
	var success bool
	var returnPosts []QuickPost
	returnPosts, success = PostCache[sub]
	return returnPosts, success
}

// AddToCacheWorker spawned to get as many reddit posts as needed
func AddToCacheWorker(sub string, wg *sync.WaitGroup, send chan<- []QuickPost) {
	defer wg.Done()
	var gottenPosts []QuickPost
	var gotPost QuickPost
	bot, _ := reddit.NewBotFromAgentFile("agent.yml", 0)
	harvest, err := bot.Listing("/r/"+sub+"/hot/", "")
	if err != nil {
		if strings.Contains(err.Error(), "400") {
			fmt.Println("Error 400. Subreddit that caused the issue:")
			fmt.Println(sub)
			return
		}
		panic(err)
	}
	for _, post := range harvest.Posts {
		mode := GuessPostType(post)
		switch {
		case mode == "link" || mode == "media":
			gotPost = QuickPost{
				Title:     post.Title,
				Score:     post.Score,
				Content:   post.URL,
				Nsfw:      post.NSFW,
				Permalink: post.Permalink,
				Sub:       sub,
			}
		case mode == "text":
			gotPost = QuickPost{
				Title:     post.Title,
				Score:     post.Score,
				Content:   post.SelfText,
				Nsfw:      post.NSFW,
				Permalink: post.Permalink,
				Sub:       sub,
			}
		}
		gottenPosts = append(gottenPosts, gotPost)
	}
	send <- gottenPosts
	if err != nil {
		panic(err)
	}
	//fmt.Println("Sent " + sub + " to cache.")
}

// PopulateCache spawns workers to add posts to the cache
func PopulateCache() {
	SubMap = subExtract("subs.json")
	fmt.Println("Populating base post cache...")
	CacheTime = time.Now().Unix() + 3600
	fmt.Println("New cache time is " + strconv.FormatInt(CacheTime, 10))
	starttime := GetMillis()
	subs := getAllSubs("subs.json")
	var wg sync.WaitGroup
	bufferSize := len(subs)
	recv := make(chan []QuickPost, bufferSize)
	for _, sub := range subs {
		wg.Add(1)
		go AddToCacheWorker(sub, &wg, recv)
	}
	wg.Wait()
	close(recv)
	for i := 0; i < bufferSize; i++ {
		var testpost QuickPost
		posts := <-recv
		for x := 0; x < len(posts); x++ {
			testpost = posts[x]
			if testpost.Sub != "" {
				break
			}
		}
		PostCache[testpost.Sub] = posts
	}
	endtime := GetMillis()
	t := endtime - starttime
	fmt.Println("Took " + strconv.FormatInt(t, 10) + "ms to add to cache.")
	CachePopulating = false
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
