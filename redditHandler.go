package main

//https://turnage.gitbooks.io/graw/content/graw.html
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

func GetFromCache(sub string) ([]QuickPost, bool) {
	var success bool
	var returnPosts []QuickPost
	returnPosts, success = PostCache[sub]
	return returnPosts, success
}

func AddToCache(sub string, addPosts []QuickPost) {
	PostCache[sub] = addPosts
}

func ClearCache() {
	for k := range PostCache {
		delete(PostCache, k)
	}
}

func PingReddit() error {
	bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
	_, err = bot.Listing("/r/all", "")
	return err
}

//gets image post
func GetMediaPost(subs []string, limit int, sort string) (QuickPost, string) {
	var returnPost QuickPost
	var cachePosts []QuickPost
	var sub string
	var success bool
	sub = subs[rand.Intn(len(subs))]
	cachePosts, success = GetFromCache(sub)
	now := time.Now().Unix()
	if now >= CacheTime {
		fmt.Println("Clearing Cache...")
		ClearCache()
		CacheTime = time.Now().Unix() + 3600
		fmt.Println("New cache time is " + strconv.FormatInt(CacheTime, 10))
	}
	switch {
	case !success:
		fmt.Println("Adding to cache.")
		var gottenPosts []QuickPost
		urlItems := []string{".jpg", ".png", ".jpeg", "gfycat", "youtube", "youtu.be", "gif", "gifv"}
		bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
		if err != nil {
			panic(err)
		}
		rand.Seed(time.Now().Unix())
		harvest, err := bot.Listing("/r/"+sub+"/"+sort, "")
		for _, post := range harvest.Posts[:limit] {
			if !strings.Contains(post.URL, "v.redd.it") && ContainsAnySubstring(post.URL, urlItems) {
				gotPost := QuickPost{
					Title:     post.Title,
					Score:     post.Score,
					Content:   post.URL,
					Nsfw:      post.NSFW,
					Permalink: post.Permalink,
				}
				gottenPosts = append(gottenPosts, gotPost)
			}
		}
		s := rand.Intn(len(gottenPosts))
		returnPost = gottenPosts[s]
		AddToCache(sub, gottenPosts)
	case success:
		fmt.Println("Found in cache.")
		s := rand.Intn(len(cachePosts))
		returnPost = cachePosts[s]
	}

	return returnPost, sub
}

func GetLinkPost(subs []string, limit int, sort string) (QuickPost, string) {
	var returnPost QuickPost
	var cachePosts []QuickPost
	var sub string
	var success bool
	sub = subs[rand.Intn(len(subs))]
	cachePosts, success = GetFromCache(sub)
	fmt.Println(success)
	if time.Now().Unix() >= CacheTime || success == false {
		fmt.Println("Adding to cache.")
		var gottenPosts []QuickPost
		bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
		if err != nil {
			panic(err)
		}
		rand.Seed(time.Now().Unix())
		harvest, err := bot.Listing("/r/"+sub+"/"+sort, "")
		for _, post := range harvest.Posts[:limit] {
			gotPost := QuickPost{
				Title:     post.Title,
				Score:     post.Score,
				Content:   post.URL,
				Nsfw:      post.NSFW,
				Permalink: post.Permalink,
			}
			gottenPosts = append(gottenPosts, gotPost)

		}
		s := rand.Intn(len(gottenPosts))
		returnPost = gottenPosts[s]
		AddToCache(sub, gottenPosts)
		CacheTime = time.Now().Unix() + 3600
	} else {
		fmt.Println("Found in cache.")
		s := rand.Intn(len(cachePosts))
		returnPost = cachePosts[s]
	}
	return returnPost, sub
}

// get self text post
func GetTextPost(subs []string, limit int, sort string) (QuickPost, string) {
	var returnPost QuickPost
	var cachePosts []QuickPost
	var sub string
	var success bool
	sub = subs[rand.Intn(len(subs))]
	cachePosts, success = GetFromCache(sub)
	fmt.Println(success)
	if time.Now().Unix() >= CacheTime || success == false {
		fmt.Println("Adding to cache.")
		var gottenPosts []QuickPost
		bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
		if err != nil {
			panic(err)
		}
		rand.Seed(time.Now().Unix())
		harvest, err := bot.Listing("/r/"+sub+"/"+sort, "")
		for _, post := range harvest.Posts[:limit] {
			gotPost := QuickPost{
				Title:     post.Title,
				Score:     post.Score,
				Content:   post.SelfText,
				Nsfw:      post.NSFW,
				Permalink: post.Permalink,
			}
			gottenPosts = append(gottenPosts, gotPost)

		}
		s := rand.Intn(len(gottenPosts))
		returnPost = gottenPosts[s]
		AddToCache(sub, gottenPosts)
		CacheTime = time.Now().Unix() + 3600
	} else {
		fmt.Println("Found in cache.")
		s := rand.Intn(len(cachePosts))
		returnPost = cachePosts[s]
	}
	return returnPost, sub
}
