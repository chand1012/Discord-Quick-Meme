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

// need tests for subreddit banning and the queue
