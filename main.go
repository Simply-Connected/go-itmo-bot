package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var config struct {
	Token string `yaml:"token"`
}

func setUpConfig() {
	r, err := os.Open("config.yaml")
	if err != nil {
		log.Panic(err)
	}
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		log.Panic(err)
	}
	r.Close()
}

func main() {
	setUpConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
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
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
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
	}
}
