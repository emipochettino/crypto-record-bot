package service

import (
	"CryptoRecordBot/internal/domain/model"
	"CryptoRecordBot/internal/domain/ports"
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Command interface {
	ShouldExecute(message telegram.Message) bool
	Execute(message telegram.Message)
}

type PriceCommand struct {
	cryptoRepository ports.CryptoRepository
	botClient        ports.BotClient
}

func (c *PriceCommand) ShouldExecute(message telegram.Message) bool {
	return message.IsCommand() && "price" == strings.ToLower(message.Command())
}

func (c *PriceCommand) Execute(message telegram.Message) {
	coinName := "bitcoin"
	if message.CommandArguments() != "" {
		coinName = message.CommandArguments()
	}
	simpleSinglePrice, err := c.cryptoRepository.GetPrice(coinName, "usd")

	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, err.Error()))
		return
	}
	price := strconv.FormatFloat(float64(simpleSinglePrice.MarketPrice), 'f', -1, 32)
	c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("%s: %s %s", coinName, simpleSinglePrice.Currency, price)))
}

func NewPriceCommand(cryptoRepository ports.CryptoRepository, botClient ports.BotClient) *PriceCommand {
	return &PriceCommand{
		cryptoRepository: cryptoRepository,
		botClient:        botClient,
	}
}

type AlertCommand struct {
	cryptoRepository ports.CryptoRepository
	alertRepository  ports.AlertRepository
	botClient        ports.BotClient
}

func (c *AlertCommand) ShouldExecute(message telegram.Message) bool {
	return message.IsCommand() && "alert" == strings.ToLower(message.Command())
}

func (c *AlertCommand) Execute(message telegram.Message) {
	arguments := strings.Split(message.CommandArguments(), " ")
	if message.CommandArguments() == "" || len(arguments) != 3 {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, "you need to indicate currency, diamond symbol and price. \n\nExample:\nBTC < 30000"))
		return
	}
	currency := arguments[0]
	diamondSymbol := arguments[1]
	priceString := arguments[2]

	if diamondSymbol != "<" && diamondSymbol != ">" {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("Argument (%s) is wrong, possible values are (<, >)", diamondSymbol)))
		return
	}
	isGreaterThan := diamondSymbol == ">"

	price, err := strconv.ParseFloat(priceString, 64)
	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("Arugment (%s) is wrong, expecting float value.\nExample 20.9113123", priceString)))
		return
	}
	coinList, err := c.cryptoRepository.GetCoinList()
	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, err.Error()))
		return
	}
	var isValidCoin bool
	for _, coin := range *coinList {
		if isValidCoin = strings.ToLower(coin.ID) == strings.ToLower(currency); isValidCoin {
			break
		}
	}
	if !isValidCoin {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("Currency %s is not valid", currency)))
		return
	}

	alert := model.MakeAlert(
		message.Chat.ID,
		message.From.ID, currency,
		isGreaterThan,
		price,
	)

	err = c.alertRepository.Create(alert)
	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, err.Error()))
		return
	}

	c.botClient.Send(telegram.NewMessage(message.Chat.ID, "Done!"))
}

func NewAlertCommand(cryptoRepository ports.CryptoRepository, botClient ports.BotClient, alertRepository ports.AlertRepository) *AlertCommand {
	return &AlertCommand{
		cryptoRepository: cryptoRepository,
		botClient:        botClient,
		alertRepository:  alertRepository,
	}
}
