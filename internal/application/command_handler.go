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

/*TODO: in order avoid mixing cqrs with ddd maybe commands can be strategies in application layer.
With that it is possible to use real domain services*/

func (c CommandHandler) Handle(message telegram.Message) {
	for _, command := range c.commands {
		if command.ShouldExecute(message) {
			command.Execute(message)
		}
	}
}
