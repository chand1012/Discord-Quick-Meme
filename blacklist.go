package main

import (
	"time"
	"fmt"
	"strconv"
)

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