package main

import (
	"gopkg.in/telegram-bot-api.v4"
)

type Bag struct {
	items []string
}

type Player struct {
	bag    *Bag
	roomID int
}

type Room struct {
	description     string
	quest           string
	items           []string
	nearRooms       map[string]int
	canGo           string
	wellcomeMessage string
	isOpen          bool
}

type Game struct {
	player Player
	rooms  []*Room
}

type GameBot struct {
	bot         *tgbotapi.BotAPI
	activeGames map[int]*Game
}

func NewGameBot() (*GameBot, error) {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		return nil, err
	}
	activeGames := make(map[int]*Game)

	return &GameBot{bot, activeGames}, nil
}
