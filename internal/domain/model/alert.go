package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Alert struct {
	ChatId        int64
	UserId        int64
	CoinName      string
	IsGreaterThan bool
	Price         float64
	CreatedAt     time.Time
}

func (a *Alert) String() string {
	diamondSymbol := "<"
	if a.IsGreaterThan {
		diamondSymbol = ">"
	}
	return fmt.Sprintf("%s %s %s", a.CoinName, diamondSymbol, a.FormattedPrice())
}

func (a *Alert) FormattedPrice() string {
	return strconv.FormatFloat(a.Price, 'f', -1, 32)
}

func MakeAlert(chatId int64, userId int64, coinName string, isGreaterThan bool, price float64) Alert {
	return Alert{
		chatId,
		userId,
		strings.ToLower(coinName),
		isGreaterThan,
		price,
		time.Now(),
	}
}
