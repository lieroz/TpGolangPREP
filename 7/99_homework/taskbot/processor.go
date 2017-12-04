package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"strings"
)

func ProcessChannelUpdate(update tgbotapi.Update) ([]*Response, error) {
	botCommand := strings.Split(update.Message.Text, " ")[0]
	for _, botComm := range botCommands {
		if botComm.FindString(botCommand) == botCommand {
			index := len(botCommand)
			if idx := strings.Index(botCommand, "_"); idx != -1 {
				index = idx
			}
			command := GetCommand(botCommand[1:index])
			responses := command().Execute(botCommand, update.Message)
			return responses, nil
		}
	}
	return nil, errNoCommand
}
