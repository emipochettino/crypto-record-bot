package domain

import "github.com/superoo7/go-gecko/v3/types"

type CryptoRepository interface {
	GetPrice(coinName string, currency string) (*types.SimpleSinglePrice, error)
}
