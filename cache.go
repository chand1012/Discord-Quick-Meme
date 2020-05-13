package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
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
	bot, err := initBot()
	if err != nil {
		fmt.Println("Error creating Bot: ", err)
		return
	}
	harvest, err := bot.Listing("/r/"+sub+"/hot/", "")
	if err != nil {
		if strings.Contains(err.Error(), "400") {
			fmt.Println("Error 400. Subreddit that caused the issue:")
			fmt.Println(sub)
			return
		}
	}
	for _, post := range harvest.Posts {
		mode := GuessPostType(post)
		if mode == "link" || mode == "media" {
			gotPost = QuickPost{
				Title:     post.Title,
				Score:     post.Score,
				Content:   post.URL,
				Nsfw:      post.NSFW,
				Permalink: post.Permalink,
				Sub:       sub,
			}
		} else {
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
}

// PopulateCache spawns workers to add posts to the cache
func PopulateCache() {
	SubMap = subExtract("subs.json")
	fmt.Println("Populating base post cache...")
	CacheTime = time.Now().Unix() + 3600
	fmt.Println("New cache time is " + strconv.FormatInt(CacheTime, 10))
	starttime := GetMillis()
	subs := getAllSubsFromMap()
	redisCommonSubs, err := getCommonSubsRedis()
	if err != nil {
		fmt.Println("Error getting common subreddits from redis: ", err)
		CachePopulating = false
		return
	}
	CommonSubsCounter = uint8(len(redisCommonSubs))
	if redisCommonSubs != nil {
		subs = append(subs, redisCommonSubs...)
	}
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
	PostCache = make(map[string][]QuickPost)
}

// This will update the counter
// if the sub is in the map, increment it
// otherwise set it to one
func updateCommonSubCounter(sub string) {
	max := uint8(60 - len(getAllSubsFromMap()))
	if CommonSubsCounter <= max && !stringInSlice(sub, getAllSubsFromMap()) {
		if GetMillis() > CommonSubsTime[sub] {
			CommonSubs[sub] = 0
			CommonSubsTime[sub] = GetMillis() + 604800000 // ms in a week
		}
		if CommonSubs[sub] < 5 {
			CommonSubs[sub]++
		}
	}
	err := saveCommonSubsRedis()
	if err != nil {
		fmt.Println("Error saving common subreddits to redis: ", err)
	}
}
