package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
	"strings"
)

type Command interface {
	Execute(string, string)
}

type GetAllTasksCommand struct {
}

func NewGetAllTasksCommand() Command {
	return &GetAllTasksCommand{}
}

func (c *GetAllTasksCommand) Execute(command, text string) {

}

type CreateTaskCommand struct {
}

func NewCreateTaskCommand() Command {
	return &CreateTaskCommand{}
}

func (c *CreateTaskCommand) Execute(command, text string) {

}

type AssignTaskCommand struct {
}

func NewAssignTaskCommand() Command {
	return &AssignTaskCommand{}
}

func (c *AssignTaskCommand) Execute(command, text string) {

}

type UnassignTaskCommand struct {
}

func NewUnassignTaskCommand() Command {
	return &UnassignTaskCommand{}
}

func (c *UnassignTaskCommand) Execute(command, text string) {

}

type ResolveTaskCommand struct {
}

func NewResolveTaskCommand() Command {
	return &ResolveTaskCommand{}
}

func (c *ResolveTaskCommand) Execute(command, text string) {

}

type GetUserAssignedTasksCommand struct {
}

func NewGetUserAssignedTasksCommand() Command {
	return &GetUserAssignedTasksCommand{}
}

func (c *GetUserAssignedTasksCommand) Execute(command, text string) {

}

type GetUserDefinedTasksCommand struct {
}

func NewGetUserDefinedTasksCommand() Command {
	return &GetUserDefinedTasksCommand{}
}

func (c *GetUserDefinedTasksCommand) Execute(command, text string) {

}

var botCommands = []*regexp.Regexp{
	regexp.MustCompile("(?i)^(/tasks)"),
	regexp.MustCompile("(?i)^(/new)"),
	regexp.MustCompile("(?i)^(/assign_\\d+)"),
	regexp.MustCompile("(?i)^(/unassign_\\d+)"),
	regexp.MustCompile("(?i)^(/resolve\\d+)"),
	regexp.MustCompile("(?i)^(/my)"),
	regexp.MustCompile("(?i)^(/own)"),
}

var commands = map[string]func() Command{
	"tasks":    NewGetAllTasksCommand,
	"new":      NewCreateTaskCommand,
	"assign":   NewAssignTaskCommand,
	"unassign": NewUnassignTaskCommand,
	"resolve":  NewResolveTaskCommand,
	"my":       NewGetUserAssignedTasksCommand,
	"own":      NewGetUserDefinedTasksCommand,
}

func getCommand(key string) func() Command {
	return commands[key]
}

func ProcessChannelUpdate(update tgbotapi.Update) error {
	botCommand := strings.Split(update.Message.Text, " ")[0]
	for _, botComm := range botCommands {
		if botComm.FindString(botCommand) == botCommand {
			botCommand = botCommand[1:]
			if idx := strings.Index(botCommand, "_"); idx != -1 {
				botCommand = botCommand[:idx]
			}
			command := getCommand(botCommand)
			command().Execute(botCommand, update.Message.Text)
			return nil
		}
	}
	return errNoCommand
}
