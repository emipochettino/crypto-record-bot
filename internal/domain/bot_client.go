package domain

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type BotClient interface {
	Send(c telegram.Chattable) (telegram.Message, error)
}
