package model

import (
	"fmt"
	"strconv"
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

func (a *Alert) String() string {
	diamondSymbol := "<"
	if a.IsGreaterThan {
		diamondSymbol = ">"
	}
	return fmt.Sprintf("%s %s %s", a.Currency, diamondSymbol, a.FormattedPrice())
}

func (a *Alert) FormattedPrice() string {
	return strconv.FormatFloat(a.Price, 'f', -1, 32)
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
