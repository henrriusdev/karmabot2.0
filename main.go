package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// global variables
var (
	httpClient *http.Client
)

// constants

func GetJson(url string, target interface{}) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func main() {
	httpClient = &http.Client{Timeout: 10 * time.Second}

	GetUpdate()
}
