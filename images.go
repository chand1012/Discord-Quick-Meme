package main

// https://github.com/vivithemage/mrisa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func imageRedditSearch(url string) string {
	var payload imagePayload
	var parsedBody returnPayload
	payload.ImageURL = url
	payload.ResizedImages = false
	var returnURL string
	var oldest int64

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Post(mrisaAddress, "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
		return "nil"
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &parsedBody)

	if err != nil {
		fmt.Println(err)
	}
	urls := parsedBody.Links
	titles := parsedBody.Descriptions
	oldest = 0
	for i := 0; i < len(urls); i++ {
		if strings.Contains(urls[i], "www.reddit.com/r/") {
			title := titles[i]
			dash := strings.Index(title, "-")
			newTitle := title[dash+2:]
			dash = strings.Index(newTitle, "-")
			timeStamp := newTitle[:dash-2]
			testTime := timeStrToSeconds(timeStamp)
			if testTime >= oldest {
				oldest = testTime
				returnURL = urls[i]
			}
		}
	}

	return returnURL

}

func imageSearch(url string) []string {
	var payload imagePayload
	var parsedBody returnPayload
	payload.ImageURL = url
	payload.ResizedImages = false

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Post(mrisaAddress, "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &parsedBody)

	if err != nil {
		fmt.Println(err)
	}
	urls := parsedBody.Links

	return urls
}
