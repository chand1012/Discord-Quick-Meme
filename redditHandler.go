package main

//https://turnage.gitbooks.io/graw/content/graw.html
import (
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

func extractInfo(subs []string, limit int) {
	url_keys := [8]string{".jpg", ".png", ".jpeg", ".gif", ".gifv", "gfycat", "youtube", "youtu.be"}
	bot, err := reddit.NewBotFromAgentFile("agent.yml", 5*time.Second)
	cfg := graw.Config{
		Subreddits: subs,
	}

}
