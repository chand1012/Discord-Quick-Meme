package main

import "testing"

//Testing
// Nothing here yet

// TestRedisConnection creates some example values and tests to see if they are
// retrieveable. Values are used in benchmarks as well
func TestRedisConnection(t *testing.T) {
	channel := "testchannel"

	err := AppendBannedSubreddit(channel, "darkjokes")

	if err != nil {
		t.Errorf("There an error appending the subreddit: %d", err)
	}

	testSlice := []string{"darkjokes"}

	returnSlice, err := GetBannedSubreddits(channel)

	if testSlice[0] != returnSlice[0] {
		t.Errorf("The subreddit expected was not returned; wanted 'darkjokes', got %s", returnSlice[0])
	}

	if err != nil {
		t.Errorf("There was an error getting the banned subreddit: %d", err)
	}

	err = UnbanSubreddit(channel, "darkjokes")

	if err != nil {
		t.Errorf("There was an error unbanning the subreddit: %d", err)
	}

	err = AppendBannedSubreddit(channel, "dankmemes")

}

// Benchmarks

// BenchmarkRedisConnection tests how long it takes to get data from the redis server.
func BenchmarkRedisConnection(b *testing.B) {

	channel := "testchannel"
	for i := 0; i < b.N; i++ {
		_, _ = GetBannedSubreddits(channel)
	}

}
