package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
)

var config struct {
	Token string
}

func setUpConfig() {
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	setUpConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err, fmt.Sprintf("\n token = %s", config.Token))
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	for update := range updates {
		go func() {
			if update.Message == nil {
				return
			}
			if !update.Message.IsCommand() {
				return
			}
			var text string
			switch update.Message.Command() {
			case "next":
				text = "next"
			case "stat":
				text = "stat"
			case "drop":
				text = "drop"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}()
	}
}
