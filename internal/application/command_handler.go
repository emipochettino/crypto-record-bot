package application

import (
	"CryptoRecordBot/internal/domain/service"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	commands []service.Command
}

func NewCommandHandler(commands ...service.Command) *CommandHandler {
	return &CommandHandler{
		commands: commands,
	}
}

func (c CommandHandler) Handle(message telegram.Message) {
	for _, command := range c.commands {
		if command.ShouldExecute(message) {
			command.Execute(message)
		}
	}
}
