package main

import (
	"CryptoRecordBot/internal/application"
	"CryptoRecordBot/internal/domain"
	"CryptoRecordBot/internal/infrastructure"
	"crypto/tls"
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	gecko "github.com/superoo7/go-gecko/v3"
)

func main() {
	geckoClient := NewGeckoClient()
	geckoRepository := infrastructure.NewGeckoRepository(geckoClient)
	token := os.Getenv("TELEGRAM_TOKEN")
	botApi, err := NewBotApi(token)
	commandHandler := application.NewCommandHandler(domain.NewGetPriceCommand(geckoRepository, botApi))
	bot, err := infrastructure.NewBot(botApi, commandHandler)
	if err != nil {
		panic(err)
	}
	bot.Start()
}

func NewGeckoClient() *gecko.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	return gecko.NewClient(httpClient)
}

func NewBotApi(token string) (*telegram.BotAPI, error) {
	if token == "" {
		return nil, fmt.Errorf("token is missing")
	}
	botApi, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	botApi.Debug = strings.EqualFold("dev", os.Getenv("PROFILE"))

	log.Printf("Authorized on account %s", botApi.Self.UserName)

	return botApi, err
}
