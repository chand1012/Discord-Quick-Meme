package main

import (
	"net/http"
	"os"
	"testing"
)

// Testing

// TestTOPGGConnection makes sure that TOPGG can be reached.
func TestTOPGGConnection(t *testing.T) {
	testBotID := "438381344943374346" // This bot's ID

	topKey := os.Getenv("TOPGG")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://top.gg/api/bots/"+testBotID, nil)
	req.Header.Set("Authorization", topKey)

	if err != nil {
		t.Errorf("There was en error setting up the HTTP request, %v", err)
	}

	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 500 {
			t.Logf("Got a status code of %d, which isn't our fault.", resp.StatusCode)
		} else {
			t.Errorf("Was expecting status code of 200, got %d", resp.StatusCode)
		}
	}

	defer resp.Body.Close()

}
