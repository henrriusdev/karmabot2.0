package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"karmabot2.0/internal/model"
)

func GetUpdatesChan(offset int64, stopChan <-chan struct{}) (UpdatesChannel, error) {
	ch := make(chan Updates, 100)
	botUrl := os.Getenv("TELEGRAM_BOT_URL")

	go func() {
		defer close(ch)
		var root Root
		for {
			select {
			case <-stopChan:
				fmt.Println("Stopping bot!")
				close(ch)
				return
			default:
			}

			if err := GetJson(botUrl, &root, offset); err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range root.Updates {
				if update.UpdateID >= offset {
					offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()
	fmt.Println("Running bot!")
	return ch, nil
}

func GetUpdates(updatesChan UpdatesChannel, stopChan <-chan struct{}, db *sql.DB) {
	plusOneRegex := regexp.MustCompile(`\+1\b`)
	minusOneRegex := regexp.MustCompile(`\-1\b`)
	botUrl := os.Getenv("TELEGRAM_BOT_URL")
	karmas := model.KarmaModel{DB: db}

	for update := range updatesChan {

		if update.Message == nil {
			continue
		}

		chat := strings.ToLower(strings.ReplaceAll(update.Message.Chat.Title, " ", "_"))

		if update.Message.Chat.IsPrivate() || update.Message.Chat.IsChannel() {
			if err := sendMessage(botUrl, update.Message.Chat.ID, "This bot can't run on private conversations and channels. Use it in a group"); err != nil {
				log.Println(err)
				return
			}
			continue
		}

		err := karmas.CreateTable(chat)
		if err != nil {
			log.Println(err)
			continue
		}

		// For bot commands
		if update.Message.IsCommand() {

			cmdText := update.Message.Command()
			switch cmdText {
			case "karma":
				userKarma, err := karmas.GetActualKarma(update.Message.From.ID, chat)
				if err != nil {
					log.Println(err)
					err = karmas.InsertUsers(update.Message.From.ID, update.Message.From.FirstName, update.Message.From.LastName, chat)
					if err != nil {
						log.Println(err)
						continue
					}

					userKarma, err = karmas.GetActualKarma(update.Message.From.ID, chat)
					if err != nil {
						log.Println(err)
						continue
					}
				}

				err = sendMessage(botUrl, update.Message.Chat.ID, update.Message.From.FirstName+" "+update.Message.From.LastName+" has "+strconv.Itoa(userKarma)+" karma points.")
				if err != nil {
					log.Println(err)
					return
				}
			case "karmalove":
				users, err := karmas.GetKarmas(chat, true)
				if err != nil {
					log.Println(err)
					continue
				}

				usersString := "Users with most karma points of " + chat + "\n"
				for i, user := range users {
					usersString += fmt.Sprintf("%d. %s has %d of karma.\n", i+1, getName(user), user.Count)
				}

				if err := sendMessage(botUrl, update.Message.Chat.ID, usersString); err != nil {
					log.Println(err)
					return
				}
			case "karmahate":
				users, err := karmas.GetKarmas(chat, false)
				if err != nil {
					log.Println(err)
					continue
				}

				usersString := "Most hated users of " + chat + "\n"
				for i, user := range users {
					usersString += fmt.Sprintf("%d. %s has %d of karma.\n", i+1, getName(user), user.Count)
				}

				if err := sendMessage(botUrl, update.Message.Chat.ID, usersString); err != nil {
					log.Println(err)
					return
				}
				continue
			}
		}

		canModify := karmas.CanModify(update.Message.From.ID, update.Message.From.FirstName, update.Message.From.LastName, chat)
		fmt.Println(canModify)
		// For +1 or -1
		if plusOneRegex.MatchString(update.Message.Text) || minusOneRegex.MatchString(update.Message.Text) {
			if update.Message.ReplyToMessage == nil {
				continue
			}

			if update.Message.From.ID == update.Message.ReplyToMessage.From.ID {
				if err := sendMessage(botUrl, update.Message.Chat.ID, "You cannot add or subtract karma yourself."); err != nil {
					log.Println(err)
					return
				}
				continue
			}

			if !canModify {
				if err := sendMessage(botUrl, update.Message.Chat.ID, "You must wait one minute to give karma."); err != nil {
					log.Println(err)
					return
				}
				continue
			} else if plusOneRegex.MatchString(update.Message.Text) {
				err = karmas.AddKarma(update.Message.From.ID, update.Message.ReplyToMessage.From.ID, update.Message.ReplyToMessage.From.FirstName, update.Message.ReplyToMessage.From.LastName, chat)
				if err != nil {
					log.Println(err)
					continue
				}

			} else if minusOneRegex.MatchString(update.Message.Text) {
				err = karmas.SubstractKarma(update.Message.From.ID, update.Message.ReplyToMessage.From.ID, update.Message.ReplyToMessage.From.FirstName, update.Message.ReplyToMessage.From.LastName, chat)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		} else {
			continue
		}

		userKarma, err := karmas.GetActualKarma(update.Message.ReplyToMessage.From.ID, chat)
		if err != nil {
			log.Println(err)
			continue
		}

		err = sendMessage(botUrl, update.Message.Chat.ID, update.Message.ReplyToMessage.From.FirstName+" has now "+strconv.Itoa(userKarma)+" karma points")
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func getName(user *model.Karma) string {
	if user.FirstName != nil && user.LastName != nil {
		return *user.FirstName + " " + *user.LastName
	} else if user.FirstName != nil && user.LastName == nil {
		return *user.FirstName
	} else {
		return "Fulanito"
	}
}
