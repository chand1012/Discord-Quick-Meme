package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type topPayload struct {
	Servers int64 `json:"server_count"`
}

//updates the server count on top.gg
func updateServerCount(uCount int64) error {
	var payload topPayload
	payload.Servers = uCount

	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://top.gg/api/bots/438381344943374346/stats", bytes.NewBuffer(data))
	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjQzODM4MTM0NDk0MzM3NDM0NiIsImJvdCI6dHJ1ZSwiaWF0IjoxNTg3MDY3OTQ5fQ.DsICD0pqFZWmR0_rSNuyxMN8b0vkFS2DH_sOfMSGkBE")
	_, err = client.Do(req)

	return err
}
