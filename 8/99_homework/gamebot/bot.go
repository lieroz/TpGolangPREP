package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	// @BotFather gives you this
	BotToken   = "424115807:AAEos5RgC4KuHRbweG7YfrrCV51W1WqsBUI"
	WebhookURL = "https://f885d09a.ngrok.io"
)

func contains(slice []string, word string) (int, bool) {
	for i, w := range slice {
		if w == word {
			return i, true
		}
	}
	return 0, false
}

func (t *GameBot) execUpdate(update tgbotapi.Update) {
	user := update.Message.From
	command := update.Message.Text

	switch {
	case strings.Contains(command, "/start"):
		t.Start(user, update)
	case strings.Contains(command, "/reset"):
		t.Reset(user, update)
	case strings.Contains(command, "осмотреться"):
		t.LookAround(user, update)
	case strings.Contains(command, "идти"):
		t.Go(command[9:], user, update)
	case strings.Contains(command, "одеть"):
		t.Dress(command[11:], user, update)
	case strings.Contains(command, "взять"):
		t.Take(command[11:], user, update)
	case strings.Contains(command, "применить"):
		t.Apply(command[19:], user, update)
	default:
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"неизвестная команда",
		))
	}
}

func MakeNewGame() *Game {
	// Новый игрок
	player := Player{nil, 0}

	// Создание комнат
	rooms := []*Room{
		{
			"ты находишься на кухне, на столе чай",
			"надо собрать рюкзак и идти в универ",
			[]string{},
			map[string]int{
				"коридор": 1,
			},
			"коридор",
			"кухня, ничего интересного",
			true,
		},
		{
			"ничего интересного",
			"",
			[]string{},
			map[string]int{
				"кухня":   0,
				"комната": 2,
				"улица":   3,
			},
			"кухня, комната, улица",
			"ничего интересного",
			true,
		},
		{
			"на столе: ключи, конспекты, на стуле - рюкзак",
			"",
			[]string{
				"ключи",
				"конспекты",
				"рюкзак",
			},
			map[string]int{
				"коридор": 1,
			},
			"коридор",
			"ты в своей комнате",
			true,
		},
		{
			"на улице уже вовсю готовятся к новому году",
			"",
			[]string{},
			map[string]int{
				"домой": 1,
			},
			"домой",
			"на улице уже вовсю готовятся к новому году",
			false,
		},
	}

	return &Game{player, rooms}
}

func startGameBot(ctx context.Context) error {
	t, err := NewGameBot()
	if err != nil {
		return err
	}

	updates := t.bot.ListenForWebhook("/")

	go http.ListenAndServe(":8081", nil)
	fmt.Println("start listen :8081")

	for update := range updates {
		t.execUpdate(update)
	}
	return nil
}

func main() {
	err := startGameBot(context.Background())
	if err != nil {
		panic(err)
	}
}
