package main

import (
	"context"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	// @BotFather gives you this
	BotToken   = "449014674:AAFuxx-aARxHr3aVC4TDemZIfh0U476m40U"
	WebhookURL = "https://ae36e2b3.ngrok.io"
)

var (
	assign   = regexp.MustCompile("/assign_\\d+")
	unassign = regexp.MustCompile("/unassign_\\d+")
	resolve  = regexp.MustCompile("/resolve_\\d+")
)

var (
	m sync.RWMutex
	// taskID int
)

func getAllTasks() string {
	if len(tasks) == 0 {
		return "Нет задач"
	}
	var response string
	formatStr := "%d. %s by @%s\n"
	for _, task := range tasks {
		response += fmt.Sprintf(formatStr, task.ID, task.Description, task.CreatedBy.UserName)
		if task.AssignedTo == nil {
			response += fmt.Sprintf("/assign_%d", task.ID)
		}
	}
	return response
}

func createNewTask(task string, user *tgbotapi.User) string {
	m.Lock()
	taskID++
	u := User{ID: user.ID, UserName: user.UserName}
	tasks = append(tasks, &Task{ID: taskID, Description: task, CreatedBy: u})
	m.Unlock()
	return fmt.Sprintf(`Задача "%s" создана, id=%d`, task, taskID)
}

func assignTask(taskID int, user *tgbotapi.User) string {
	for _, task := range tasks {
		if task.ID == taskID {
			response := fmt.Sprintf(`Задача "%s" назначена на вас`, task.Description)
			if task.AssignedTo != nil && task.AssignedTo.ID != user.ID {
				response += fmt.Sprintf(` %d:Задача "%s" назначена на @%s`, task.AssignedTo.ID, task.Description, user.UserName)
			}
			task.AssignedTo = &User{ID: user.ID, UserName: user.UserName}
			return response
		}
	}
	return "Нет задачи с таким идентификатором!"
}

func processMessage(message *tgbotapi.Message) string {
	command := strings.Split(message.Text, " ")[0]
	var response string
	switch command {
	case "/tasks":
		response = getAllTasks()

	case "/new":
		task := strings.TrimLeft(message.Text, command)
		response = createNewTask(task[1:], message.From)

	case assign.FindString(command):
		id, _ := strconv.Atoi(strings.TrimLeft(command, "/assign_"))
		response = assignTask(id, message.From)

	case unassign.FindString(command):
		fmt.Println("UNASSIGN")

	case resolve.FindString(command):
		fmt.Println("RESOLVE")

	case "/my":
		fmt.Println("MY")

	case "/own":
		fmt.Println("OWN")
	}
	return response
}

func showAllTasks() error {
	return nil
}

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return err
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		return err
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":8080", nil)
	fmt.Println("server started on port 8080")

	for update := range updates {
		result := processMessage(update.Message)
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			result,
		))
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
