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

//ResetBlacklist just guess what this does
func ResetBlacklist() {
	fmt.Println("Resetting Blacklist...")
	Blacklist = make(map[string][]QuickPost)
	BlacklistTime = time.Now().Unix() + (3600 * 3)
	fmt.Println("New Blacklist time is " + strconv.FormatInt(BlacklistTime, 10))
}

// UpdateBlacklistTime updates clears the blacklist and updates the time
func UpdateBlacklistTime() {
	nowTime := time.Now().Unix()
	if nowTime >= BlacklistTime {
		ResetBlacklist()
	}
}

//CheckBlacklist compares the blacklist to the given post
func CheckBlacklist(channel string, post QuickPost) bool {
	count := 0
	cacheLength := len(Blacklist[channel])
	for _, cachedPost := range Blacklist[channel] {
		count++
		if count >= cacheLength {
			ResetBlacklist()
			return false
		}
		if post == cachedPost {
			return true
		}
	}
	return false
}

//AddToBlacklist add post to blacklist
func AddToBlacklist(channel string, post QuickPost) {
	Blacklist[channel] = append(Blacklist[channel], post)
}

// GetFromCache pulls post from the cache
func GetFromCache(sub string) ([]QuickPost, bool) {
	var success bool
	var returnPosts []QuickPost
	returnPosts, success = PostCache[sub]
	return returnPosts, success
}

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
	fmt.Println("Populating base post cache...")
	CacheTime = time.Now().Unix() + 3600
	fmt.Println("New cache time is " + strconv.FormatInt(CacheTime, 10))
	starttime := GetMillis()
	subs := []string{"dankmemes", "funny", "memes", "comedyheaven", "blackpeopletwitter", "whitepeopletwitter", "MemeEconomy", "therewasanattempt", "wholesomememes", "instant_regret", "ahegao", "Artistic_Hentai", "Hentai", "MonsterGirl", "slimegirls", "wholesomehentai", "quick_hentai", "HentaiParadise", "jokes", "darkjokes", "antijokes", "fiftyfifty", "UpliftingNews", "news", "worldnews", "FloridaMan", "nottheonion"}
	var wg sync.WaitGroup
	bufferSize := len(subs)
	recv := make(chan []QuickPost, bufferSize)
	for _, sub := range subs {
		wg.Add(1)
		go AddToCacheWorker(sub, &wg, recv)
	}
	/*
		for _, sub := range CommonSubs {
			wg.Add(1)
			go AddToCacheWorker(sub, &wg, recv)
		}
	*/
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
	return avg / 2
}
