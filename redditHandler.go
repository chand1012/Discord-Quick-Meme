package main

import (
	"fmt"
	"math/rand"
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
func AddToCacheWorker(sub string, wg *sync.WaitGroup, bot reddit.Scanner, send chan<- []QuickPost, mode string) {
	defer wg.Done()
	var gottenPosts []QuickPost
	var gotPost QuickPost
	harvest, err := bot.Listing("/r/"+sub+"/hot/", "")
	for _, post := range harvest.Posts {
		switch mode {
		case "link":
			gotPost = QuickPost{
				Title:     post.Title,
				Score:     post.Score,
				Content:   post.URL,
				Nsfw:      post.NSFW,
				Permalink: post.Permalink,
				Sub:       sub,
			}
		case "text":
			gotPost = QuickPost{
				Title:     post.Title,
				Score:     post.Score,
				Content:   post.SelfText,
				Nsfw:      post.NSFW,
				Permalink: post.Permalink,
				Sub:       sub,
			}
		case "media":
			urlItems := []string{".jpg", ".png", ".jpeg", "gfycat", "youtube", "youtu.be", "gif", "gifv"}
			if !strings.Contains(post.URL, "v.redd.it") && ContainsAnySubstring(post.URL, urlItems) {
				gotPost = QuickPost{
					Title:     post.Title,
					Score:     post.Score,
					Content:   post.URL,
					Nsfw:      post.NSFW,
					Permalink: post.Permalink,
					Sub:       sub,
				}
			}
		}
		gottenPosts = append(gottenPosts, gotPost)
	}
	send <- gottenPosts
	if err != nil {
		panic(err)
	}
}

// PopulateCache spawns workers to add posts to the cache
func PopulateCache() {
	fmt.Println("Populating base post cache...")
	CacheTime = time.Now().Unix() + 3600
	fmt.Println("New cache time is " + strconv.FormatInt(CacheTime, 10))
	starttime := GetMillis()
	mediaSubs := []string{"dankmemes", "funny", "memes", "comedyheaven", "blackpeopletwitter", "whitepeopletwitter", "MemeEconomy", "therewasanattempt", "wholesomememes", "instant_regret", "ahegao", "Artistic_Hentai", "Hentai", "MonsterGirl", "slimegirls", "wholesomehentai", "quick_hentai", "HentaiParadise"}
	textSubs := []string{"jokes", "darkjokes", "antijokes"}
	linkSubs := []string{"fiftyfifty", "UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"}
	var wg sync.WaitGroup
	bufferSize := len(mediaSubs) + len(textSubs) + len(linkSubs)
	recv := make(chan []QuickPost, bufferSize)
	bot, _ := reddit.NewBotFromAgentFile("agent.yml", 0)
	for _, sub := range mediaSubs {
		wg.Add(1)
		go AddToCacheWorker(sub, &wg, bot, recv, "media")
	}
	for _, sub := range textSubs {
		wg.Add(1)
		go AddToCacheWorker(sub, &wg, bot, recv, "text")
	}
	for _, sub := range linkSubs {
		wg.Add(1)
		go AddToCacheWorker(sub, &wg, bot, recv, "link")
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

	var subList []string
	var urlItems []string
	var sub string

	var success bool

	var s int

	subList = []string{"dankmemes", "funny", "memes", "comedyheaven", "blackpeopletwitter", "whitepeopletwitter", "MemeEconomy", "therewasanattempt", "wholesomememes", "instant_regret", "jokes", "darkjokes", "antijokes", "UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion", "fiftyfifty"}
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
		} else if ContainsAnySubstring(sub, subList) {
			s = rand.Intn(gottenLength)
			returnPost = gottenPosts[s]
			go PopulateCache()
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
	return returnPost, sub
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
	fmt.Println(avg / 2)
	return avg / 2
}
