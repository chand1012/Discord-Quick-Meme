package main

import (
	"database/sql"
	"testing"
)

// Tests

// TestChannelDB tests the channel database operations
func TestChannelDB(t *testing.T) {
	testChannel := "0000000000"
	testNSFW := true
	name := "testChannel"
	guildID := "111111111111111"
	err := AddChannelToDB(testChannel, testNSFW, name, guildID)

	if err != nil && err != sql.ErrNoRows {
		t.Errorf("There was an error adding to the channel DB, %v", err)
	}

	returnNSFW, returnName, guild, err := GetChannelFromDB(testChannel)

	if err != nil {
		t.Errorf("There was an error getting channel from the DB, %v", err)
	}

	if !returnNSFW || returnName != "testChannel" || guild != guildID {
		t.Errorf("Expected set values, got %t and %s", returnNSFW, returnName)
	}

	err = RemoveChannelFromDB(testChannel)

	if err != nil {
		t.Errorf("There was an error removing the channel from the database: %v", err)
	}

}

func TestCharsetRecovery(t *testing.T) {
	err := FixDatabaseTableCharset()

	if err != nil {
		t.Errorf("There was an error setting the charset of the tables: %v", err)
	}
}

func TestSubRedditBanning(t *testing.T) {
	testChannel := "0000000000"
	testSub := "imgoingtohellforthis" // this sub died, gonna use it for testing

	err := SetBannedSubreddit(testChannel, testSub)

	if err != nil {
		t.Errorf("There was an error setting banned subreddit: %v", err)
	}

	bannedSubs, err := GetAllBannedSubs(testChannel)

	if err != nil {
		t.Errorf("There was an error getting banned subreddit: %v", err)
	}

	if bannedSubs[0] != testSub {
		t.Errorf("Unexpected values, expected %s but got %s", testSub, bannedSubs[0])
	}

	err = RemoveBannedSubreddit(testChannel, testSub)

	if err != nil {
		t.Errorf("There was an error removing banned subreddit: %v", err)
	}
}

// need tests for the queue
