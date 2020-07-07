package bot

import (
	"fmt"
	"github.com/Skazzi00/go-itmo-bot/storage"
	"github.com/yanzay/tbot/v2"
	"math"
)

const (
	messagePrefix = "/"
	messageSuffix = "$"
	registerCmd   = messagePrefix + "register" + messageSuffix
	statCmd       = messagePrefix + "stat" + messageSuffix
	dropCmd       = messagePrefix + "drop" + messageSuffix
	nextCmd       = messagePrefix + "next" + messageSuffix
)

type Params struct {
	Logger       tbot.Logger
	RoomsStorage storage.RoomsStorage
	Token        string
}

func NewBot(params *Params) *tbot.Server {
	server := tbot.New(params.Token, tbot.WithLogger(params.Logger))
	client := server.Client()
	roomsStorage := params.RoomsStorage
	logger := params.Logger
	server.HandleMessage(registerCmd, func(message *tbot.Message) {
		userName := message.From.Username
		chatID := message.Chat.ID
		_, ok := roomsStorage.GetUserPoints(chatID, userName)
		if ok {
			client.SendMessage(chatID, "You are already registered", tbot.OptReplyToMessageID(message.MessageID))
		} else {
			err := roomsStorage.SetUserPoints(chatID, userName, 0)
			if err != nil {
				logger.Error(err)
				return
			}
			client.SendMessage(chatID, "You are registered", tbot.OptReplyToMessageID(message.MessageID))
		}
	})
	server.HandleMessage(dropCmd, func(message *tbot.Message) {
		userName := message.From.Username
		chatID := message.Chat.ID
		_, ok := roomsStorage.GetUserPoints(chatID, userName)
		if !ok {
			client.SendMessage(chatID, "You are not registered", tbot.OptReplyToMessageID(message.MessageID))
			return
		}
		err := roomsStorage.IncUserPoints(chatID, userName)
		if err != nil {
			logger.Error(err)
			return
		}
		client.SendMessage(chatID, "Dropped!", tbot.OptReplyToMessageID(message.MessageID))
	})
	server.HandleMessage(nextCmd, func(message *tbot.Message) {
		chatID := message.Chat.ID
		room, err := roomsStorage.GetRoom(chatID)
		if err != nil {
			logger.Error(err)
			return
		}
		if len(room) == 0 {
			client.SendMessage(chatID, "Nobody")
			return
		}
		var result string
		cur := math.MaxInt32
		for k, v := range room {
			if v < cur {
				cur, result = v, k
			}
		}
		client.SendMessage(chatID, fmt.Sprintf("@%v is next", result))
	})
	server.HandleMessage(statCmd, func(message *tbot.Message) {
		chatID := message.Chat.ID
		room, err := roomsStorage.GetRoom(chatID)
		if err != nil {
			logger.Error(err)
			return
		}
		if len(room) == 0 {
			client.SendMessage(chatID, "Nobody is registered")
			return
		}
		var result string
		for k, v := range room {
			result += fmt.Sprintf("%v: %v\n", k, v)
		}
		client.SendMessage(chatID, fmt.Sprintf("Scoreboard\n%v", result))
	})
	return server
}
