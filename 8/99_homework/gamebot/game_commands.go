package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"strings"
)

func (t *GameBot) Start(user *tgbotapi.User, update tgbotapi.Update) {
	t.activeGames[user.ID] = MakeNewGame()
	t.bot.Send(tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"добро пожаловать в игру!",
	))
}

func (t *GameBot) Reset(user *tgbotapi.User, update tgbotapi.Update) {
	delete(t.activeGames, user.ID)
	game := MakeNewGame()
	t.activeGames[user.ID] = game
	t.bot.Send(tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"состояние игры сброшено",
	))
}

func (t *GameBot) LookAround(user *tgbotapi.User, update tgbotapi.Update) {
	game := t.activeGames[user.ID]
	room := game.rooms[game.player.roomID]
	answer := room.description

	if game.player.bag != nil {
		if strings.Contains(room.quest, "собрать рюкзак и ") {
			room.quest = strings.Replace(room.quest, "собрать рюкзак и ", "", -1)
		} else if strings.Contains(room.quest, "собрать рюкзак и ") {
			room.quest = strings.Replace(room.quest, "собрать рюкзак", "", -1)
		}
	}

	if room.quest != "" {
		answer += ", " + room.quest
	}
	answer += "."
	if len(room.nearRooms) > 0 {
		answer += " можно пройти - " + room.canGo
		answer += "."
	}
	t.bot.Send(tgbotapi.NewMessage(
		update.Message.Chat.ID,
		answer,
	))
}

func (t *GameBot) Go(roomName string, user *tgbotapi.User, update tgbotapi.Update) {
	game := t.activeGames[user.ID]
	room := game.rooms[game.player.roomID]
	roomID, ok := room.nearRooms[roomName]

	if !ok {
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"нет пути в "+roomName,
		))
		return
	}

	newRoom := game.rooms[roomID]
	if newRoom.isOpen {
		game.player.roomID = roomID
		answer := newRoom.wellcomeMessage + "."
		if len(newRoom.nearRooms) > 0 {
			answer += " можно пройти - " + newRoom.canGo
			answer += "."
		}
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			answer,
		))
	} else {
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"дверь закрыта",
		))
	}
}

func (t *GameBot) Dress(itemName string, user *tgbotapi.User, update tgbotapi.Update) {
	game := t.activeGames[user.ID]
	room := game.rooms[game.player.roomID]
	if itemName == "рюкзак" {
		i, ok := contains(room.items, itemName)
		if ok {
			game.player.bag = new(Bag)

		}
		room.items = append(room.items[:i], room.items[i+1:]...)
		room.description = strings.Replace(room.description, ", на стуле - рюкзак", "", -1)
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"вы одели: "+itemName,
		))
	}
}

func (t *GameBot) Take(itemName string, user *tgbotapi.User, update tgbotapi.Update) {
	game := t.activeGames[user.ID]
	room := game.rooms[game.player.roomID]

	if game.player.bag != nil {
		i, ok := contains(room.items, itemName)
		if ok {
			game.player.bag.items = append(game.player.bag.items, itemName)
		} else {
			t.bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"нет такого",
			))
			return
		}
		room.items = append(room.items[:i], room.items[i+1:]...)

		lastIndex := strings.LastIndex(room.description, itemName) + len(itemName)
		if lastIndex < len(room.description) && room.description[lastIndex:lastIndex+1] == "," {
			room.description = strings.Replace(room.description, " "+itemName+",", "", -1)
		} else {
			room.description = strings.Replace(room.description, " "+itemName, "", -1)
		}

		if len(room.items) == 0 {
			room.description = "пустая комната"
		}

		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"предмет добавлен в инвентарь: "+itemName,
		))
	} else {
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"некуда класть",
		))
	}
}

func (t *GameBot) Apply(args string, user *tgbotapi.User, update tgbotapi.Update) {
	game := t.activeGames[user.ID]
	room := game.rooms[game.player.roomID]

	ar := strings.Split(args, " ")

	ok := false
	if game.player.bag != nil {
		_, ok = contains(game.player.bag.items, ar[0])
	}
	if !ok {
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"нет предмета в инвентаре - "+ar[0],
		))
		return
	}

	if ar[1] != "дверь" {
		t.bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"не к чему применить",
		))
		return
	}
	if args == "ключи дверь" {
		// Ищем комнату с дверью
		haveDoor := false
		var closedRoomID int
		for _, roomID := range room.nearRooms {
			if !game.rooms[roomID].isOpen {
				haveDoor = true
				closedRoomID = roomID
				break
			}
		}
		if haveDoor {
			game.rooms[closedRoomID].isOpen = true

			t.bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"дверь открыта",
			))
		} else {
			t.bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"здесь нет дверей",
			))
		}
	}
}
