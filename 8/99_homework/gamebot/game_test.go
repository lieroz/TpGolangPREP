package main

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"sync/atomic"
	// "encoding/json"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func init() {
	// upd global var for testing
	// we use patched version of gopkg.in/telegram-bot-api.v4 ( WebhookURL const -> var)
	WebhookURL = "http://127.0.0.1:8081"
	BotToken = "_golangcourse_test_gamebot"
}

var (
	client = &http.Client{Timeout: time.Second}
)

// TDS is Telegram Dummy Server
type TDS struct {
	*sync.Mutex
	Answers map[int]string
}

func NewTDS() *TDS {
	return &TDS{
		Mutex:   &sync.Mutex{},
		Answers: make(map[int]string),
	}
}

func (srv *TDS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/getMe", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true,"result":{"id":` +
			strconv.Itoa(BotChatID) +
			`,"is_bot":true,"first_name":"game_test_bot","username":"game_test_bot"}}`))
	})
	mux.HandleFunc("/setWebhook", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true,"result":true,"description":"Webhook was set"}`))
	})
	mux.HandleFunc("/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		chatID, _ := strconv.Atoi(r.FormValue("chat_id"))
		text := r.FormValue("text")
		srv.Lock()
		srv.Answers[chatID] = text
		srv.Unlock()

		//fmt.Println("TDS sendMessage", chatID, text)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("unknown command %s", r.URL.Path))
	})

	handler := http.StripPrefix("/bot"+BotToken, mux)
	handler.ServeHTTP(w, r)
}

const (
	Alice int = 256
	Bob   int = 512

	BotChatID = 100501
)

var (
	users = map[int]*tgbotapi.User{
		Alice: &tgbotapi.User{
			ID:           Alice,
			FirstName:    "Alice",
			LastName:     "null",
			UserName:     "alice",
			LanguageCode: "ru",
			IsBot:        false,
		},
		Bob: &tgbotapi.User{
			ID:           Bob,
			FirstName:    "Bob",
			LastName:     "undef",
			UserName:     "bob",
			LanguageCode: "ru",
			IsBot:        false,
		},
	}

	updID uint64
	msgID uint64
)

func SendMsgToBot(userID int, text string) error {
	// reqText := `{
	// 	"update_id":175894614,
	// 	"message":{
	// 		"message_id":29,
	// 		"from":{"id":133250764,"is_bot":false,"first_name":"Vasily Romanov","username":"rvasily","language_code":"ru"},
	// 		"chat":{"id":133250764,"first_name":"Vasily Romanov","username":"rvasily","type":"private"},
	// 		"date":1512168732,
	// 		"text":"THIS SEND FROM USER"
	// 	}
	// }`

	atomic.AddUint64(&updID, 1)
	myUpdID := atomic.LoadUint64(&updID)

	// better have it per user, but lazy now
	atomic.AddUint64(&msgID, 1)
	myMsgID := atomic.LoadUint64(&msgID)

	user, ok := users[userID]
	if !ok {
		return fmt.Errorf("no user %s", userID)
	}

	upd := &tgbotapi.Update{
		UpdateID: int(myUpdID),
		Message: &tgbotapi.Message{
			MessageID: int(myMsgID),
			From:      user,
			Chat: &tgbotapi.Chat{
				ID:        int64(user.ID),
				FirstName: user.FirstName,
				UserName:  user.UserName,
				Type:      "private",
			},
			Text: text,
			Date: int(time.Now().Unix()),
		},
	}
	reqData, _ := json.Marshal(upd)

	reqBody := bytes.NewBuffer(reqData)
	req, _ := http.NewRequest(http.MethodPost, WebhookURL, reqBody)
	_, err := client.Do(req)
	return err
}

type testCase struct {
	user    int
	command string
	answers map[int]string
}

type singlePlayerCase struct {
	command string
	answer  string
}

var singlePlayerCases = [][]singlePlayerCase{

	[]singlePlayerCase{
		{
			"/start",
			"добро пожаловать в игру!",
		},
		{ // действие осмотреться
			"осмотреться",
			"ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор.",
		},
		{ // действие идти
			"идти коридор",
			"ничего интересного. можно пройти - кухня, комната, улица.",
		},
		{
			"идти комната",
			"ты в своей комнате. можно пройти - коридор.",
		},
		{
			"осмотреться",
			"на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор.",
		},
		{ // действие одеть
			"одеть рюкзак",
			"вы одели: рюкзак",
		},
		{ // действие взять
			"взять ключи",
			"предмет добавлен в инвентарь: ключи",
		},
		{
			"взять конспекты",
			"предмет добавлен в инвентарь: конспекты",
		},
		{
			"идти коридор",
			"ничего интересного. можно пройти - кухня, комната, улица.",
		},
		{ // действие применить
			"применить ключи дверь",
			"дверь открыта",
		},
		{
			"идти улица",
			"на улице уже вовсю готовятся к новому году. можно пройти - домой.",
		},
		{
			"/reset",
			"состояние игры сброшено",
		},
	},
	[]singlePlayerCase{
		{
			"/start",
			"добро пожаловать в игру!",
		},
		{
			"осмотреться",
			"ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор.",
		},
		{ // придёт топать в универ голодным :(
			"завтракать",
			"неизвестная команда",
		},
		{ // через стены ходить нельзя
			"идти комната",
			"нет пути в комната",
		},
		{
			"идти коридор",
			"ничего интересного. можно пройти - кухня, комната, улица.",
		},
		{
			"применить ключи дверь",
			"нет предмета в инвентаре - ключи",
		},
		{
			"идти комната",
			"ты в своей комнате. можно пройти - коридор.",
		},
		{
			"осмотреться",
			"на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор.",
		},
		{ // надо взять рюкзак сначала
			"взять ключи",
			"некуда класть",
		},
		{
			"одеть рюкзак",
			"вы одели: рюкзак",
		},
		{ // состояние изменилось
			"осмотреться",
			"на столе: ключи, конспекты. можно пройти - коридор.",
		},
		{
			"взять ключи",
			"предмет добавлен в инвентарь: ключи",
		},
		{ // неизвестный предмет
			"взять телефон",
			"нет такого",
		},
		{ // предмента уже нет в комнате - мы его взяли
			"взять ключи",
			"нет такого",
		},
		{ // состояние изменилось
			"осмотреться",
			"на столе: конспекты. можно пройти - коридор.",
		},
		{
			"взять конспекты",
			"предмет добавлен в инвентарь: конспекты",
		},
		{ // состояние изменилось
			"осмотреться",
			"пустая комната. можно пройти - коридор.",
		},
		{
			"идти коридор",
			"ничего интересного. можно пройти - кухня, комната, улица.",
		},
		{
			"идти кухня",
			"кухня, ничего интересного. можно пройти - коридор.",
		},
		{ // состояние изменилось
			"осмотреться",
			"ты находишься на кухне, на столе чай, надо идти в универ. можно пройти - коридор.",
		},
		{
			"идти коридор",
			"ничего интересного. можно пройти - кухня, комната, улица.",
		},
		{ //условие не удовлетворено
			"идти улица",
			"дверь закрыта",
		},
		{ //состояние изменилось
			"применить ключи дверь",
			"дверь открыта",
		},
		{ // нет предмета
			"применить телефон шкаф",
			"нет предмета в инвентаре - телефон",
		},
		{ // предмет есть, но применить его к этому нельзя
			"применить ключи шкаф",
			"не к чему применить",
		},
		{
			"идти улица",
			"на улице уже вовсю готовятся к новому году. можно пройти - домой.",
		},
		{
			"/reset",
			"состояние игры сброшено",
		},
	},
}

func TestSinglePlayer(t *testing.T) {

	tds := NewTDS()
	ts := httptest.NewServer(tds)
	tgbotapi.APIEndpoint = ts.URL + "/bot%s/%s"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := startGameBot(ctx)
		if err != nil {
			t.Fatalf("startGameBot error: %s", err)
		}
	}()
	
	// give server time to start
    time.Sleep(10 * time.Millisecond)

	for _, gameSession := range singlePlayerCases {

		for idx, item := range gameSession {
			tds.Lock()
			tds.Answers = make(map[int]string)
			tds.Unlock()

			caseName := fmt.Sprintf("[case%d, %d: %s]", idx, Alice, item.command)
			err := SendMsgToBot(Alice, item.command)
			if err != nil {
				t.Fatalf("% SendMsgToBot error: %s", caseName, err)
			}
			// give TDS time to process request
			time.Sleep(10 * time.Millisecond)

			expected := map[int]string{
				Alice: item.answer,
			}

			tds.Lock()
			result := reflect.DeepEqual(tds.Answers, expected)
			if !result {
				t.Fatalf("%s bad results:\n\tWant: %v\n\tHave: %v", caseName, expected, tds.Answers)
			}
			tds.Unlock()
		}
	}
}
