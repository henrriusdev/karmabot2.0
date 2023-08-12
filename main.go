package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// global variables
var (
	httpClient *http.Client
	collection *mongo.Collection
)

// constants
const (
	connectionString = "mongodb://localhost:27017/"
	dbName           = "gokarmas"
	colName          = "karmas"
)

func init() {
	// client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongonnected")

	collection = client.Database(dbName).Collection(colName)

	// collection instance
	fmt.Println("Mongollection is ready")
}

func GetUpdate() {
	url := "https://api.telegram.org/bot5917051686:AAEf5hPR7stgvKb2Ig38IXfydEk88vpnUnI/getUpdates"

	var update Updates

	if err := GetJson(url, &update); err != nil {
		fmt.Printf("error getting update: %s\n", err.Error())
		return
	} else {
		fmt.Printf("Update received!")
	}
}

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
