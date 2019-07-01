package main

//https://turnage.gitbooks.io/graw/content/graw.html
import (
	"math/rand"
	"time"

	"github.com/turnage/graw/reddit"
)

//gets image post
func GetPost(subs []string, limit int) (int32, string, string, bool) {
	var scores []int32
	var urls []string
	var titles []string
	var nsfws []bool
	bot, err := reddit.NewBotFromAgentFile("agent.yml", 5*time.Second)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().Unix())
	sub := subs[rand.Intn(len(subs))]
	harvest, err := bot.Listing("/r/"+sub, "")

	for _, post := range harvest.Posts[:limit] {
		scores = append(scores, post.Score)
		urls = append(urls, post.URL)
		titles = append(titles, post.Title)
		nsfws = append(nsfws, post.NSFW)
	}

	s := rand.Intn(limit)

	return scores[s], urls[s], titles[s], nsfws[s]

}
