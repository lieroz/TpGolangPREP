package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	commandMutex sync.RWMutex
)

type Command interface {
	Execute(string, *tgbotapi.Message) []*Response
}

type GetAllTasksCommand struct {
}

func NewGetAllTasksCommand() Command {
	return &GetAllTasksCommand{}
}

func (c *GetAllTasksCommand) Execute(command string, message *tgbotapi.Message) []*Response {
	var responses []*Response

	if len(tasks) == 0 {
		responses = append(responses, &Response{ChatID: message.Chat.ID, Message: errNoTasks.Error()})
		return responses
	}

	var msg string
	for i, task := range tasks {
		if i > 0 {
			msg += "\n\n"
		}

		msg += fmt.Sprintf(taskFormatResponse+"\n", task.ID, task.Description, task.CreatedBy.UserName)

		if task.AssignedTo == nil {
			msg += fmt.Sprintf(assignFormat, task.ID)
		} else {
			if task.AssignedTo.ID == message.From.ID {
				formatStr := assigneeMe + "\n" + unassignFormat + " " + resolveFormat
				msg += fmt.Sprintf(formatStr, task.ID, task.ID)
			} else {
				msg += fmt.Sprintf(assigneeUser, task.AssignedTo.UserName)
			}
		}
	}

	responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})
	return responses
}

type CreateTaskCommand struct {
}

func NewCreateTaskCommand() Command {
	return &CreateTaskCommand{}
}

func (c *CreateTaskCommand) Execute(command string, message *tgbotapi.Message) []*Response {
	var responses []*Response
	task := strings.TrimLeft(message.Text, command)[1:]
	AddTask(task, message.From)
	msg := fmt.Sprintf(taskCreatedResponse, task, taskID)
	responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})
	return responses
}

type AssignTaskCommand struct {
}

func NewAssignTaskCommand() Command {
	return &AssignTaskCommand{}
}

func (c *AssignTaskCommand) Execute(command string, message *tgbotapi.Message) []*Response {
	var responses []*Response
	id, _ := strconv.Atoi(strings.TrimLeft(command, "/assign_"))

	for _, task := range tasks {
		if task.ID == id {
			msg := fmt.Sprintf(taskAssignedToYouResponse, task.Description)
			responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})

			if task.AssignedTo == nil && task.CreatedBy.ID != message.From.ID {
				msg = fmt.Sprintf(taskAssignedToUserResponse, task.Description, message.From.UserName)
				responses = append(responses, &Response{ChatID: int64(task.CreatedBy.ID), Message: msg})
			}
			if task.AssignedTo != nil && task.AssignedTo.ID != message.From.ID {
				msg = fmt.Sprintf(taskAssignedToUserResponse, task.Description, message.From.UserName)
				responses = append(responses, &Response{ChatID: int64(task.AssignedTo.ID), Message: msg})
			}

			task.AssignedTo = message.From
			return responses
		}
	}

	return nil
}

type UnassignTaskCommand struct {
}

func NewUnassignTaskCommand() Command {
	return &UnassignTaskCommand{}
}

func (c *UnassignTaskCommand) Execute(command string, message *tgbotapi.Message) []*Response {
	var responses []*Response
	id, _ := strconv.Atoi(strings.TrimLeft(command, "/unassign_"))

	for _, task := range tasks {
		if task.ID == id {

			if task.AssignedTo.ID != message.From.ID {
				msg := fmt.Sprintf(errTaskNotYour.Error())
				responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})
			} else {
				msg := fmt.Sprintf(taskUnassignAcceptedResponse)
				responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})

				msg = fmt.Sprintf(taskWithoutImplementerResponse, task.Description)
				responses = append(responses, &Response{ChatID: int64(task.CreatedBy.ID), Message: msg})

				task.AssignedTo = nil
			}

			return responses
		}
	}

	return nil
}

type ResolveTaskCommand struct {
}

func NewResolveTaskCommand() Command {
	return &ResolveTaskCommand{}
}

func (c *ResolveTaskCommand) Execute(command string, message *tgbotapi.Message) []*Response {
	var responses []*Response
	id, _ := strconv.Atoi(strings.TrimLeft(command, "/resolve_"))

	for i, task := range tasks {
		if task.ID == id {

			if task.AssignedTo.ID != message.From.ID {
				msg := fmt.Sprintf(errTaskNotYour.Error())
				responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})
			} else {
				msg := fmt.Sprintf(taskDoneResponse, task.Description)
				responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})

				msg = fmt.Sprintf(taskDoneByResponse, task.Description, message.From.UserName)
				responses = append(responses, &Response{ChatID: int64(task.CreatedBy.ID), Message: msg})

				tasks = append(tasks[:i], tasks[i+1:]...)
			}

			return responses
		}
	}

	return nil
}

type GetUserAssignedTasksCommand struct {
}

func NewGetUserAssignedTasksCommand() Command {
	return &GetUserAssignedTasksCommand{}
}

func (c *GetUserAssignedTasksCommand) Execute(command string, message *tgbotapi.Message) []*Response {
	var responses []*Response
	var msg string
	counter := 0
	for _, task := range tasks {

		if task.AssignedTo != nil && task.AssignedTo.ID == message.From.ID {
			counter++
			if counter > 1 {
				msg += "\n\n"
			}

			msg += fmt.Sprintf(taskFormatResponse+"\n", task.ID, task.Description, task.CreatedBy.UserName)
			formatStr := unassignFormat + " " + resolveFormat
			msg += fmt.Sprintf(formatStr, task.ID, task.ID)
		}
	}

	responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})
	return responses
}

type GetUserDefinedTasksCommand struct {
}

func NewGetUserDefinedTasksCommand() Command {
	return &GetUserDefinedTasksCommand{}
}

func (c *GetUserDefinedTasksCommand) Execute(command string, message *tgbotapi.Message) []*Response {
	var responses []*Response
	var msg string
	counter := 0
	for _, task := range tasks {

		if task.CreatedBy.ID == message.From.ID {
			counter++
			if counter > 1 {
				msg += "\n\n"
			}

			msg += fmt.Sprintf(taskFormatResponse+"\n", task.ID, task.Description, task.CreatedBy.UserName)

			if task.AssignedTo == nil || task.AssignedTo.ID != task.CreatedBy.ID {
				msg += fmt.Sprintf(assignFormat, task.ID)
			} else {
				formatStr := unassignFormat + " " + resolveFormat
				msg += fmt.Sprintf(formatStr, task.ID, task.ID)
			}
		}
	}

	responses = append(responses, &Response{ChatID: message.Chat.ID, Message: msg})
	return responses
}

var botCommands = []*regexp.Regexp{
	regexp.MustCompile("(?i)^(/tasks)"),
	regexp.MustCompile("(?i)^(/new)"),
	regexp.MustCompile("(?i)^(/assign_\\d+)"),
	regexp.MustCompile("(?i)^(/unassign_\\d+)"),
	regexp.MustCompile("(?i)^(/resolve_\\d+)"),
	regexp.MustCompile("(?i)^(/my)"),
	regexp.MustCompile("(?i)^(/owner)"),
}

var commands = map[string]func() Command{
	"tasks":    NewGetAllTasksCommand,
	"new":      NewCreateTaskCommand,
	"assign":   NewAssignTaskCommand,
	"unassign": NewUnassignTaskCommand,
	"resolve":  NewResolveTaskCommand,
	"my":       NewGetUserAssignedTasksCommand,
	"owner":    NewGetUserDefinedTasksCommand,
}

func GetCommand(key string) func() Command {
	commandMutex.RLock()
	result := commands[key]
	commandMutex.RUnlock()
	return result
}
