package service

import (
	"CryptoRecordBot/internal/domain/ports"
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

type AlertService struct {
	cryptoRepository ports.CryptoRepository
	alertRepository  ports.AlertRepository
	botClient        ports.BotClient
}

func NewAlertService(cryptoRepository ports.CryptoRepository, alertRepository ports.AlertRepository, botClient ports.BotClient) *AlertService {
	return &AlertService{
		cryptoRepository: cryptoRepository,
		alertRepository:  alertRepository,
		botClient:        botClient,
	}
}

func (s *AlertService) AlertByCoinName() {
	coinNames, err := s.alertRepository.FindCoinNames()
	if err != nil {
		log.Println("Something went wrong trying to find coinNames", err)
		return
	}
	for _, coinName := range coinNames {
		price, err := s.cryptoRepository.GetPrice(coinName, "usd")
		if err != nil {
			log.Printf("Something went wrong trying to find the price of %s. Error: %s\n", coinName, err.Error())
			return
		}

		alerts, err := s.alertRepository.FindByCoinName(coinName)
		if err != nil {
			log.Println("Something went wrong trying to find alerts by coinNames", err)
			return
		}

		for _, alert := range alerts {
			if alert.IsGreaterThan && alert.Price < float64(price.MarketPrice) {
				s.botClient.Send(telegram.NewMessage(
					alert.ChatId,
					fmt.Sprintf("%s has a price of %s and it is higher than %s",
						alert.CoinName,
						strconv.FormatFloat(float64(price.MarketPrice), 'f', -1, 32),
						alert.FormattedPrice())),
				)
				s.alertRepository.Delete(alert)
			} else if !alert.IsGreaterThan && alert.Price > float64(price.MarketPrice) {
				s.botClient.Send(telegram.NewMessage(
					alert.ChatId,
					fmt.Sprintf("%s has a price of %s and it is lower than %s",
						alert.CoinName,
						strconv.FormatFloat(float64(price.MarketPrice), 'f', -1, 32),
						alert.FormattedPrice())),
				)
				s.alertRepository.Delete(alert)
			}
		}
	}
}
