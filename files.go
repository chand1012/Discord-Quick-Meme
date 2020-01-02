package main

// https://github.com/vivithemage/mrisa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type imagePayload struct {
	ImageURL      string `json:"image_url"`
	ResizedImages bool   `json:"resized_images"`
}

type returnPayload struct {
	Links         []string `json:"links"`
	Descriptions  []string `json:"descriptions"`
	SimilarImages []string `json:"similar_images"`
	Titles        []string `json:"titles"`
}

func imageSearch(url string) returnPayload {
	var payload imagePayload
	var parsedBody returnPayload
	payload.ImageURL = url
	payload.ResizedImages = false

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}
	// this will be changed to a different server that will be hidden in
	// data.json
	resp, err := http.Post("http://localhost:5000/search", "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &parsedBody)

	if err != nil {
		fmt.Println(err)
	}

	return parsedBody

}
