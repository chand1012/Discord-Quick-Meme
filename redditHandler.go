package main

//https://turnage.gitbooks.io/graw/content/graw.html
import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
)

//gets image post
func GetMediaPost(subs []string, limit int) (int32, string, string, bool, string, string) {
	var scores []int32
	var urls []string
	var titles []string
	var nsfws []bool
	var links []string
	urlItems := []string{".jpg", ".png", ".jpeg", "gfycat", "youtube", "youtu.be", "gif", "gifv"}
	bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().Unix())
	sub := subs[rand.Intn(len(subs))]
	harvest, err := bot.Listing("/r/"+sub, "")
	fmt.Println("Getting posts.....")
	for _, post := range harvest.Posts[:limit] {
		if !strings.Contains(post.URL, "v.redd.it") && ContainsAnySubstring(post.URL, urlItems) {
			scores = append(scores, post.Score)
			urls = append(urls, post.URL)
			titles = append(titles, post.Title)
			nsfws = append(nsfws, post.NSFW)
			links = append(links, post.Permalink)
		}
	}

	s := rand.Intn(len(urls))

	return scores[s], urls[s], titles[s], nsfws[s], links[s], sub

}

// get self text post
func GetTextPost(subs []string, limit int) (int32, string, string, bool, string, string) {
	var scores []int32
	var text []string
	var titles []string
	var nsfws []bool
	var links []string
	bot, err := reddit.NewBotFromAgentFile("agent.yml", 0)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().Unix())
	sub := subs[rand.Intn(len(subs))]
	harvest, err := bot.Listing("/r/"+sub, "")

	for _, post := range harvest.Posts[:limit] {
		scores = append(scores, post.Score)
		text = append(text, post.SelfText)
		titles = append(titles, post.Title)
		nsfws = append(nsfws, post.NSFW)
		links = append(links, post.Permalink)
	}

	s := rand.Intn(len(text))

	return scores[s], text[s], titles[s], nsfws[s], links[s], sub
}
