package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"sync"
)

type Task struct {
	ID          int
	Description string
	CreatedBy   *tgbotapi.User
	AssignedTo  *tgbotapi.User
}

var (
	taskID int
	mutex  sync.Mutex

	users []*tgbotapi.User
	tasks []*Task
)

func AddUser(user *tgbotapi.User) {
	users = append(users, user)
}

func AddTask(descr string, user *tgbotapi.User) {
	m.Lock()
	taskID++
	tasks = append(tasks, &Task{ID: taskID, Description: descr, CreatedBy: user})
	m.Unlock()
}
