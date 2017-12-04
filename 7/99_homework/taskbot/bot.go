package main

import (
	"context"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
)

var (
	// @BotFather gives you this
	BotToken   = "449014674:AAFuxx-aARxHr3aVC4TDemZIfh0U476m40U"
	WebhookURL = "https://7644ddf1.ngrok.io"
)

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
		responses, err := ProcessChannelUpdate(update)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				err.Error(),
			))
		}
		for _, r := range responses {
			bot.Send(tgbotapi.NewMessage(
				r.ChatID,
				r.Message,
			))
		}
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
