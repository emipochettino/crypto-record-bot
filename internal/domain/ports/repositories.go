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
	FindByChatIDAndUserID(chatID int64, userID int64) ([]model.Alert, error)
	Create(alert model.Alert) error
	Delete(alert model.Alert) (bool, error)
	FindCoinNames() ([]string, error)
	FindByCoinName(coinName string) ([]model.Alert, error)
}
