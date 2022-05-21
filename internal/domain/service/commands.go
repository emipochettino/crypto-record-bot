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
	//TODO: extract this common functionality
	price := strconv.FormatFloat(float64(simpleSinglePrice.MarketPrice), 'f', -1, 32)
	c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("%s: %s %s", coinName, simpleSinglePrice.Currency, price)))
}

func NewPriceCommand(cryptoRepository ports.CryptoRepository, botClient ports.BotClient) *PriceCommand {
	return &PriceCommand{
		cryptoRepository: cryptoRepository,
		botClient:        botClient,
	}
}

type CreateAlertCommand struct {
	cryptoRepository ports.CryptoRepository
	alertRepository  ports.AlertRepository
	botClient        ports.BotClient
}

func (c *CreateAlertCommand) ShouldExecute(message telegram.Message) bool {
	return message.IsCommand() && "createalert" == strings.ToLower(message.Command())
}

func (c *CreateAlertCommand) Execute(message telegram.Message) {
	arguments := strings.Split(message.CommandArguments(), " ")
	if message.CommandArguments() == "" || len(arguments) != 3 {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, "you need to indicate coin name, diamond symbol and price. \n\nExample:\nbitcoin < 30000"))
		return
	}
	coinName := arguments[0]
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
	coins, err := c.cryptoRepository.GetCoinList()
	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, err.Error()))
		return
	}
	var isValidCoin bool
	for _, coin := range *coins {
		if isValidCoin = strings.ToLower(coin.ID) == strings.ToLower(coinName); isValidCoin {
			break
		}
	}
	if !isValidCoin {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("CoinName %s is not valid", coinName)))
		return
	}

	alert := model.MakeAlert(
		message.Chat.ID,
		message.From.ID,
		coinName,
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

func NewCreateAlertCommand(cryptoRepository ports.CryptoRepository, botClient ports.BotClient, alertRepository ports.AlertRepository) *CreateAlertCommand {
	return &CreateAlertCommand{
		cryptoRepository: cryptoRepository,
		botClient:        botClient,
		alertRepository:  alertRepository,
	}
}

type DeleteAlertCommand struct {
	alertRepository ports.AlertRepository
	botClient       ports.BotClient
}

func (c *DeleteAlertCommand) ShouldExecute(message telegram.Message) bool {
	return message.IsCommand() && "deletealert" == strings.ToLower(message.Command())
}

func (c *DeleteAlertCommand) Execute(message telegram.Message) {
	coinName := message.CommandArguments()
	if coinName == "" {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, "you need to indicate coinName.\nExample: bitcoin"))
		return
	}

	alert := model.MakeAlert(
		message.Chat.ID,
		message.From.ID,
		coinName,
		false,
		0,
	)

	isDeleted, err := c.alertRepository.Delete(alert)
	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, err.Error()))
		return
	}

	if isDeleted {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, "Done!"))
		return
	} else {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf("Alert with coinName (%s) not found!", coinName)))
		return
	}
}

func NewDeleteAlertCommand(botClient ports.BotClient, alertRepository ports.AlertRepository) *DeleteAlertCommand {
	return &DeleteAlertCommand{
		botClient:       botClient,
		alertRepository: alertRepository,
	}
}

type ListAlertsCommand struct {
	alertRepository ports.AlertRepository
	botClient       ports.BotClient
}

func (c *ListAlertsCommand) ShouldExecute(message telegram.Message) bool {
	return message.IsCommand() && "listalerts" == strings.ToLower(message.Command())
}

func (c *ListAlertsCommand) Execute(message telegram.Message) {
	alerts, err := c.alertRepository.FindByChatIDAndUserID(message.Chat.ID, message.From.ID)
	if err != nil {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, err.Error()))
		return
	}
	if len(alerts) == 0 {
		c.botClient.Send(telegram.NewMessage(message.Chat.ID, "you dont have alerts "))
		return
	}

	var alertsString []string
	for _, alert := range alerts {
		alertsString = append(alertsString, alert.String())
	}
	c.botClient.Send(telegram.NewMessage(message.Chat.ID, strings.Join(alertsString, "\n")))

}

func NewListAlertsCommand(botClient ports.BotClient, alertRepository ports.AlertRepository) *ListAlertsCommand {
	return &ListAlertsCommand{
		botClient:       botClient,
		alertRepository: alertRepository,
	}
}
