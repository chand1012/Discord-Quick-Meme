package jsonHandler

import (
	"encoding/json"
	"os"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Keys struct {
	botId	     string	 `json:token`
	clientId     string  `json:client_id`
	clientSecret string  `json:client_secret`
	userAgent	 string  `json:user_agent`
}

func extract(filename string) {
	jsonfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonfile.Close()

	rawjson, _ := ioutil.ReadAll(jsonFile)

	var keys Keys

	json.Unmarshall(rawjson, &keys)

	return keys, err
}