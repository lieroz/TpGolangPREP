package main

import (
	"context"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	// @BotFather gives you this
	BotToken   string // "449014674:AAFuxx-aARxHr3aVC4TDemZIfh0U476m40U"
	WebhookURL string // "https://7644ddf1.ngrok.io"
	Port       string
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

	server := &http.Server{Addr: Port, Handler: nil}
	go server.ListenAndServe()
	fmt.Println("server started on port " + Port)

	c, finish := context.WithCancel(ctx)
	go func() {
		select {
		case <-c.Done():
			server.Shutdown(ctx)
			os.Exit(0)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for {
			s := <-sigChan
			switch s {
			case syscall.SIGINT:
				finish()
			}
		}
	}()

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

func getEnvVars() {
	BotToken = os.Getenv("TOKEN")
	WebhookURL = os.Getenv("URL")
	Port = os.Getenv("PORT")
}

func main() {
	getEnvVars()

	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
