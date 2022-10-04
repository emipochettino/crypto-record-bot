package bootstrap

import (
	"CryptoRecordBot/internal/application"
	"CryptoRecordBot/internal/domain/service"
	"CryptoRecordBot/internal/infrastructure"
	"CryptoRecordBot/internal/infrastructure/client"
	"CryptoRecordBot/internal/infrastructure/persistance"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gecko "github.com/superoo7/go-gecko/v3"
)

type App struct {
	Bot *infrastructure.Bot
}

func NewApp() *App {
	whiteListString := os.Getenv("WHITE_LIST")
	var whiteList []int64
	if strings.TrimSpace(whiteListString) != "" {
		whiteListStrSlice := strings.Split(whiteListString, ",")
		whiteList = make([]int64, len(whiteListStrSlice))
		for i, s := range whiteListStrSlice {
			whiteList[i], _ = strconv.ParseInt(s, 10, 64)
		}
	}
	geckoClient := NewGeckoClient()
	geckoRepository := client.NewGeckoRepository(geckoClient)
	token := os.Getenv("TELEGRAM_TOKEN")
	if len(token) == 0 {
		panic("TELEGRAM_TOKEN is not set")
	}
	botApi, err := NewBotApi(token)
	db := persistance.NewDB()
	alertRepository := persistance.NewAlertRepository(db)
	commandHandler := application.NewCommandHandler(
		service.NewPriceCommand(geckoRepository, botApi),
		service.NewCreateAlertCommand(geckoRepository, botApi, alertRepository),
		service.NewDeleteAlertCommand(botApi, alertRepository),
		service.NewListAlertsCommand(botApi, alertRepository),
	)
	bot, err := infrastructure.NewBot(botApi, commandHandler, whiteList)
	if err != nil {
		panic(err)
	}

	alertService := service.NewAlertService(geckoRepository, alertRepository, botApi)
	go func() {
		for {
			alertService.AlertByCoinName()
			time.Sleep(3 * time.Minute)
		}
	}()

	return &App{
		Bot: bot,
	}
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
