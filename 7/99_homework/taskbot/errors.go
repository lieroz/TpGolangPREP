package main

import (
	"errors"
)

var (
	errNoCommand = errors.New("Нет такой команды")
	errNoTask    = errors.New("Нет задачи с таким идентификатором")

	errNoTasks     = errors.New("Нет задач")
	errTaskNotYour = errors.New("Задача не на вас")
)
