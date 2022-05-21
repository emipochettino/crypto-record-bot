package ports

import (
	"CryptoRecordBot/internal/domain/model"
	"github.com/superoo7/go-gecko/v3/types"
)

type CryptoRepository interface {
	GetPrice(coinName string, currency string) (*types.SimpleSinglePrice, error)
	GetCoinList() (*types.CoinList, error)
}

type AlertRepository interface {
	FindByID(chatID int64, userID int64, currency string, isGreaterThan bool) *model.Alert
	FindByChatIDAndUserID(chatID int64, userID int64) ([]model.Alert, error)
	Create(alert model.Alert) error
	Delete(alert model.Alert) error
}
