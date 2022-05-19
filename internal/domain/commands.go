package domain

import (
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Command interface {
	ShouldExecute(message telegram.Message) bool
	Execute(message telegram.Message) telegram.MessageConfig
}

type GetPriceCommand struct {
	cryptoRepository CryptoRepository
	botClient        BotClient
}

func (c *GetPriceCommand) ShouldExecute(message telegram.Message) bool {
	return message.IsCommand() && "price" == strings.ToLower(message.Command())
}

func (c *GetPriceCommand) Execute(message telegram.Message) telegram.MessageConfig {
	coinName := "bitcoin"
	if message.CommandArguments() != "" {
		coinName = message.CommandArguments()
	}
	simpleSinglePrice, err := c.cryptoRepository.GetPrice(coinName, "usd")

	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, err.Error()))
		return telegram.NewMessage(message.Chat.ID, err.Error())
	}
	price := strconv.FormatFloat(float64(simpleSinglePrice.MarketPrice), 'f', -1, 32)
	c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("%s: %s %s", coinName, simpleSinglePrice.Currency, price)))

	return telegram.NewMessage(message.Chat.ID, fmt.Sprintf("%s: %s %s", coinName, simpleSinglePrice.Currency, price))
}

func NewGetPriceCommand(cryptoRepository CryptoRepository, botClient BotClient) *GetPriceCommand {
	return &GetPriceCommand{
		cryptoRepository: cryptoRepository,
		botClient:        botClient,
	}
}
