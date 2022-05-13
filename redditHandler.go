package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
)

// Replaces the "reddit.NewBotFromAgentFile" with a simple function call. Uses
// getRedditEnv and gets data from environment
func initBot() (reddit.Bot, error) {
	var agent agentFile = getRedditEnv()

	app := reddit.App{
		ID:     agent.ClientID,
		Secret: agent.ClientSecret,
	}

	bot, err := reddit.NewBot(
		reddit.BotConfig{
			Agent: agent.UserAgent,
			App:   app,
			Rate:  0,
		},
	)
	return bot, err
}

// GuessPostType get the type of the post so a
func GuessPostType(post *reddit.Post) string {
	selfText := post.SelfText
	urlContent := post.URL
	if selfText == "" {
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
	bot, err := initBot()

	if err != nil {
		log.Println("Error pinging Reddit:", err)
		return err
	}

	_, err = bot.Listing("/r/all", "")
	if err != nil {
		log.Println("Error pinging Reddit:", err)
	}
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
	cachePosts, success = GetFromCache(sub)
	now := time.Now().Unix()
	if now >= CacheTime && !CachePopulating {
		log.Println("Clearing Cache...")
		ClearCache()
		success = false
		CacheTime = time.Now().Unix() + 3600
		log.Println("New cache time is " + strconv.FormatInt(CacheTime, 10))
		CachePopulating = true
		go PopulateCache()
	}
	if !success {
		starttime := GetMillis()
		bot, err := initBot()
		if err != nil {
			log.Println("Error creating new Reddit bot:", err)
			return QuickPost{}, ""
		}
		rand.Seed(time.Now().Unix())
		log.Println("Adding r/" + sub + " to cache.")
		harvest, err := bot.Listing("/r/"+sub+"/"+sort, "") // the bot is locking up here
		if err != nil {
			log.Println("Error getting posts from sub '", sub, "':", err)
			return QuickPost{}, sub
		}
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
			log.Println("Nothing to cache! Discarding....")
		} else if ContainsAnySubstring(sub, subList) {
			s = rand.Intn(gottenLength)
			returnPost = gottenPosts[s]
		} else {
			AddToCache(sub, gottenPosts)
			CacheTime = time.Now().Unix() + 1800
			s = rand.Intn(gottenLength)
			returnPost = gottenPosts[s]
			endtime := GetMillis()
			t := endtime - starttime
			log.Println("Took " + strconv.FormatInt(t, 10) + "ms to add to cache.")
		}
	} else {
		log.Println("Found r/" + sub + " in cache.")
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
	var n int = len(posts)
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
		allSubs = append(allSubs, value...)
	}
	return allSubs
}
