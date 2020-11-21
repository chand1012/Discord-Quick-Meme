package main

import (
	"testing"
)

func TestCache(t *testing.T) {
	PostCache = make(map[string][]QuickPost)
	CachePopulating = true
	CacheTime = 0

	PopulateCache()

	if CachePopulating == true {
		t.Error("Cache says still populating, even though it should be finished.")
	}

	if CacheTime == 0 {
		t.Error("CacheTime is still zero.")
	}

	if len(PostCache) == 0 {
		t.Error("Cache is empty, did not populate!")
	}

	ClearCache()

	if len(PostCache) != 0 {
		t.Error("Cache is not empty, did not clear!")
	}

}
