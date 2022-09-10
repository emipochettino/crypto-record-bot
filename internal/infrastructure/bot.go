package infrastructure

import (
	"CryptoRecordBot/internal/application"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	BotApi         *telegram.BotAPI
	commandHandler *application.CommandHandler
	whiteList      []int64
}

func (bot *Bot) sendMessage(messageConfig telegram.MessageConfig) {
	bot.BotApi.Send(messageConfig)
}

func NewBot(botApi *telegram.BotAPI, commandHandler *application.CommandHandler, whiteList []int64) (*Bot, error) {
	return &Bot{BotApi: botApi, commandHandler: commandHandler, whiteList: whiteList}, nil
}

func (bot *Bot) Start() {
	updateConfig := telegram.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.BotApi.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			if len(bot.whiteList) > 0 {
				if !contains(bot.whiteList, update.Message.From.ID) {
					log.Printf("User with id %d and username %s is trying to use the bot", update.Message.From.ID, update.Message.From.UserName)
					continue
				}
			}
			go bot.commandHandler.Handle(*update.Message)
		}
	}
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
