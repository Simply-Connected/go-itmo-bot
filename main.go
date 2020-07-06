package main

import (
	"github.com/Skazzi00/go-itmo-bot/bot"
	"github.com/Skazzi00/go-itmo-bot/storage/pudge"
	"github.com/yanzay/tbot/v2"
	"log"
	"os"
)

func setupParams() *bot.Params {
	return &bot.Params{
		Logger:       tbot.BasicLogger{},
		RoomsStorage: pudge.NewRoomsStorage(),
		Token:        os.Getenv("BOT_TOKEN"),
	}
}

func main() {
	params := setupParams()
	garbageBot := bot.NewBot(params)
	err := garbageBot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
