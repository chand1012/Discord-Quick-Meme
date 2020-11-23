package main

import (
	"testing"
	"time"
)

func TestBlacklist(t *testing.T) {
	testChannel := "test" + randString(14)
	testPost := QuickPost{
		Title:     "Test",
		Score:     400,
		Content:   "https://github.com/chand1012/Discord-Quick-Meme",
		Nsfw:      false,
		Permalink: "https://github.com/chand1012/Discord-Quick-Meme",
		Sub:       "github",
	}

	ResetBlacklist()
	testTime := time.Now().Unix()

	if testTime >= BlacklistTime {
		t.Errorf("Blacklist time is less than system time! BlacklistTime: %d; System Time: %d", BlacklistTime, testTime)
	}

	if len(Blacklist) != 0 {
		t.Errorf("Blacklist is not empty! Got %d.", len(Blacklist))
	}

	Blacklist[testChannel] = append(Blacklist[testChannel], testPost)

	if len(Blacklist) != 1 {
		t.Errorf("Blacklist does not have 1 element! Expected 1, got %d.", len(Blacklist))
	}

	if len(Blacklist[testChannel]) != 1 {
		t.Errorf("Blacklist for test channel '%s' does not have 1 element! Expected 1, got %d.", testChannel, len(Blacklist[testChannel]))
	}

	blacklistCheck := CheckBlacklist(testChannel, testPost)

	if !blacklistCheck {
		t.Errorf("Blacklist does not contain the test post!")
	}

	BlacklistTime = 0

	UpdateBlacklistTime()

	if len(Blacklist) != 0 {
		t.Errorf("Blacklist was not reset! Got %d.", len(Blacklist))
	}

}
