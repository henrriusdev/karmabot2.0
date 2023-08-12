package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	//"strings"
)

func GetUpdate() {
	plusOneRegex := regexp.MustCompile(`\+1\b`)
	minusOneRegex := regexp.MustCompile(`\-1\b`)
	url := os.Getenv("TELEGRAM_BOT_URL")

	var updates Updates

	if err := GetJson(url, &updates); err != nil {
		log.Fatal(err)
	}

	for _, update := range updates.Result {
		// chat := strings.ToLower(strings.ReplaceAll(update.Message.Chat.Title, " ", "_"))

		if update.Message.Chat.IsPrivate() || update.Message.Chat.IsChannel() {
			msg := "This bot can't run on DM and channels, use it in a group or supergroup"

			fmt.Print(msg)
			continue
		}

		if update.Message.IsCommand() {
			cmdText := update.Message.Command()
			switch cmdText {
			case "karma":
				fmt.Println("Get the karma")
			case "karmalove":
				fmt.Println("Get the users with most karma")
			case "karmahate":
				fmt.Println("Get the users with minor quantity of karma")
			}
		}

		if plusOneRegex.MatchString(update.Message.Text) || minusOneRegex.MatchString(update.Message.Text) {
			if true {
				fmt.Println("This if check the time updated")
			} else if plusOneRegex.MatchString(update.Message.Text) {
				fmt.Println("plus one to the karma user")
			} else if minusOneRegex.MatchString(update.Message.Text) {
				fmt.Println("substract one to the karma user")
			}
		} else {
			continue
		}

		// sends the message of new karma
	}
}
