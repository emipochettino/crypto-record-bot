package main

import (
	"CryptoRecordBot/internal/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	app.Bot.Start()
}
