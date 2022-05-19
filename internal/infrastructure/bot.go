package infrastructure

import (
	"CryptoRecordBot/internal/application"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotApi         *telegram.BotAPI
	commandHandler *application.CommandHandler
}

func (bot *Bot) sendMessage(messageConfig telegram.MessageConfig) {
	bot.BotApi.Send(messageConfig)
}

func NewBot(botApi *telegram.BotAPI, commandHandler *application.CommandHandler) (*Bot, error) {
	return &Bot{BotApi: botApi, commandHandler: commandHandler}, nil
}

func (bot *Bot) Start() {
	updateConfig := telegram.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.BotApi.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			go bot.commandHandler.Handle(*update.Message)
		}
	}
}
