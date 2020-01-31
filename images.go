package main

// https://github.com/vivithemage/mrisa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

// this if for "how long ago"
// ex: "5 hours ago"
// ex: "4 days ago"
func timeStrToSeconds(stamp string) int64 {
	var stampLength string
	var finalTime int64
	lengths := []string{"second", "minute", "hour", "day", "year"}
	for _, l := range lengths {
		if strings.Contains(stamp, l) {
			stampLength = l
			break
		}
	}
	spaceIndex := strings.Index(stamp, " ")
	timeUnknown, _ := strconv.ParseInt(stamp[:spaceIndex], 10, 64)
	switch stampLength {
	case "second":
		finalTime = timeUnknown

	case "minute":
		finalTime = timeUnknown * 60

	case "hour":
		finalTime = timeUnknown * 3600

	case "day":
		finalTime = timeUnknown * 3600 * 24

	case "year":
		finalTime = timeUnknown * 3600 * 8760
	}
	return finalTime
}

func imageSearch(url string) (string, string) {
	var payload imagePayload
	var parsedBody returnPayload
	payload.ImageURL = url
	payload.ResizedImages = false
	var returnURL string
	var returnTitle string
	var oldest int64

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
	}
	// this will be changed to a different server that will be hidden in
	// data.json
	resp, err := http.Post("http://127.0.0.1:5000/search", "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
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
			fmt.Println(testTime)
			if testTime >= oldest {
				oldest = testTime
				returnURL = urls[i]
				returnTitle = newTitle
			}
		}
	}

	return returnURL, returnTitle

}
