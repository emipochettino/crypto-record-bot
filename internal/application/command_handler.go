package application

import (
	"CryptoRecordBot/internal/domain"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	commands []domain.Command
}

func NewCommandHandler(commands ...domain.Command) *CommandHandler {
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
