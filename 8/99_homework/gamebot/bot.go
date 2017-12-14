package main

import (
	"context"
)

var (
	// @BotFather gives you this
	BotToken   = "XXX"
	WebhookURL = "https://525f2cb5.ngrok.io"
)

/*
func startGameBot(ctx context.Context) error {
	// сюда пишите ваш код
	return nil
}
*/

func main() {
	err := startGameBot(context.Background())
	if err != nil {
		panic(err)
	}
}
