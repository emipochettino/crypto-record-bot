package infrastructure

import (
	coingecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type GeckoRepository struct {
	geckoClient *coingecko.Client
}

func (g GeckoRepository) GetPrice(coinName string, currency string) (*types.SimpleSinglePrice, error) {
	return g.geckoClient.SimpleSinglePrice(coinName, currency)
}

func NewGeckoRepository(geckoClient *coingecko.Client) *GeckoRepository {
	return &GeckoRepository{
		geckoClient: geckoClient,
	}
}
