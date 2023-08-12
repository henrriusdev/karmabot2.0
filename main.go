package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

// global variables
var (
	httpClient *http.Client
)

// constants

func GetJson(baseUrl string, target interface{}, offset int64) error {
	data := url.Values{}
	data.Add("allowed_updates", "['message']")
	data.Add("offset", strconv.Itoa(int(offset)))

	completeUrl := baseUrl + "/getUpdates" + "?" + data.Encode()
	resp, err := httpClient.Get(completeUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func sendMessage(baseUrl string, chatID int64, message string) error {
	data := url.Values{}
	data.Add("chat_id", strconv.FormatInt(chatID, 10))
	data.Add("text", message)

	completeUrl := baseUrl + "/sendMessage" + "?" + data.Encode()

	resp, err := httpClient.Get(completeUrl)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()
	return nil
}

func main() {
	httpClient = &http.Client{Timeout: 10 * time.Second}

	godotenv.Load()

	stopChan := make(chan struct{})

	db, err := openDB(os.Getenv("KARMA_CONN_STRING"))
	if err != nil {
		log.Println(err)
		return
	}

	updates, err := GetUpdatesChan(0, stopChan)
	if err != nil {
		log.Fatal(err)
	}

	GetUpdates(updates, stopChan, db)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
