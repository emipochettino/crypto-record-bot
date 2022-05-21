package client

import (
	"CryptoRecordBot/internal/domain/ports"
	coingecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type GeckoRepository struct {
	geckoClient *coingecko.Client
}

func (g GeckoRepository) GetPrice(coinName string, currency string) (*types.SimpleSinglePrice, error) {
	return g.geckoClient.SimpleSinglePrice(coinName, currency)
}

func (g GeckoRepository) GetCoinList() (*types.CoinList, error) {
	return g.geckoClient.CoinsList()
}

func NewGeckoRepository(geckoClient *coingecko.Client) ports.CryptoRepository {
	return &GeckoRepository{
		geckoClient: geckoClient,
	}
}
