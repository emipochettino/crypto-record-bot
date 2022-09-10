package client

import (
	"CryptoRecordBot/internal/domain/model"
	"CryptoRecordBot/internal/domain/ports"
	"encoding/json"
	"fmt"
	coingecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type GeckoRepository struct {
	geckoClient *coingecko.Client
}

func (g GeckoRepository) GetPrice(coinName string, currency string) (*types.SimpleSinglePrice, error) {
	return g.geckoClient.SimpleSinglePrice(coinName, currency)
}

func (g GeckoRepository) GetPriceWith24hsChange(coinName string) (*model.SimplePrice, error) {
	response, err := g.geckoClient.MakeReq(fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd&include_24hr_change=true",
		coinName,
	))
	mappedResponse := make(map[string]model.SimplePrice)
	err = json.Unmarshal(response, &mappedResponse)
	if err != nil {
		return nil, err
	}

	if len(mappedResponse) == 0 {
		return nil, fmt.Errorf("token %s not found", coinName)
	}

	return &model.SimplePrice{
		Id: coinName,
		PriceWithChange: model.PriceWithChange{
			Usd:          mappedResponse[coinName].Usd,
			Usd24HChange: mappedResponse[coinName].Usd24HChange,
		},
	}, nil
}

func (g GeckoRepository) GetCoinList() (*types.CoinList, error) {
	return g.geckoClient.CoinsList()
}

func NewGeckoRepository(geckoClient *coingecko.Client) ports.CryptoRepository {
	return &GeckoRepository{
		geckoClient: geckoClient,
	}
}
