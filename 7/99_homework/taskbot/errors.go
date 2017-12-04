package main

import (
	"errors"
)

var (
	ErrNoCommand = errors.New("нет такой команды...")
	ErrNoTask    = errors.New("нет задачи с таким идентификатором...")
)
