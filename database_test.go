package main

import (
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

// Tests

// TestChannelDB tests the channel database operations
func TestChannelDB(t *testing.T) {
	godotenv.Load()
	testChannel := "0000000000"
	testNSFW := true
	name := "testChannel"

	err := AddChannelToDB(testChannel, testNSFW, name)

	if err != nil && err != mongo.ErrNoDocuments {
		t.Errorf("There was an error adding to the channel DB, %v", err)
	}

	returnNSFW, returnName, err := GetChannelFromDB(testChannel)

	if err != nil {
		t.Errorf("There was an error getting channel from the DB, %v", err)
	}

	if !returnNSFW && returnName != "testChannel" {
		t.Errorf("Expected set values, got %t and %s", returnNSFW, returnName)
	}

	err = RemoveChannelFromDB(testChannel)

	if err != nil {
		t.Errorf("There was an error removing the channel from the database: %v", err)
	}

}

// need tests for subreddit banning and the queue
