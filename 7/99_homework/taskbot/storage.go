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
	taskID    int
	taskMutex sync.Mutex

	tasks []*Task
)

func AddTask(descr string, user *tgbotapi.User) {
	taskMutex.Lock()
	taskID++
	tasks = append(tasks, &Task{ID: taskID, Description: descr, CreatedBy: user})
	taskMutex.Unlock()
}
