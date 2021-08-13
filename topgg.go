package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type topPayload struct {
	Servers int64 `json:"server_count"`
}

//updates the server count on top.gg
func updateServerCount(uCount int64, topKey string) {
	var payload topPayload
	payload.Servers = uCount

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("There was an error while encoding JSON:", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://top.gg/api/bots/"+botID+"/stats", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("There was an error while creating the request:", err)
		return
	}

	req.Header.Set("Authorization", topKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		fmt.Println(req)
		fmt.Println(resp)
		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan() && i < 10; i++ {
			fmt.Println(scanner.Text())
		}
	}

	if err != nil {
		fmt.Println("There was an error while setting the server count:", err)
	}

	defer resp.Body.Close()

}
