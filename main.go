package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/recoilme/pudge"
	"io/ioutil"
	"log"
)

type Config struct {
	Token string
}
type UserName string

func getConfig() Config {
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panic(err)
	}
	var config Config
	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Panic(err)
	}
	return config
}

func main() {
	config := getConfig()
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
	rooms, err := pudge.Open("rooms", &pudge.Config{SyncInterval: 1})
	if err != nil {
		log.Panic(err)
	}
	defer rooms.Close()

	type Room map[UserName]int
	for update := range updates {
		go func() {

			if update.Message == nil {
				return
			}
			if !update.Message.IsCommand() {
				return
			}
			var room Room
			_ = rooms.Get(update.Message.Chat.ID, &room)
			user := update.Message.From
			var result string
			switch update.Message.Command() {
			case "next":
				cur := room[UserName(user.UserName)]
				result = user.UserName
				for k, v := range room {
					if v < cur {
						cur, result = v, string(k)
					}
				}
			case "stat":
				if len(room) == 0 {
					room[UserName(user.UserName)] = 0
				}
				for k, v := range room {
					result += fmt.Sprintf("%v: %v\n", k, v)
				}
			case "drop":
				room[UserName(user.UserName)] += 1
				result = "Dropped!"
			}
			rooms.Set(update.Message.Chat.ID, &room)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}()
	}
}
