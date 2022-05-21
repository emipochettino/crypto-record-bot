package model

import (
	"time"
)

type Alert struct {
	ChatId        int64
	UserId        int64
	Currency      string
	IsGreaterThan bool
	Price         float64
	CreatedAt     time.Time
}

func MakeAlert(chatId int64, userId int64, currency string, isGreaterThan bool, price float64) Alert {
	return Alert{
		chatId,
		userId,
		currency,
		isGreaterThan,
		price,
		time.Now(),
	}
}
