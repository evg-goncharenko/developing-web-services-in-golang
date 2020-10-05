package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
находясь в корне репы
git subtree push --prefix 4/5_bot heroku master
*/

const (
	BotToken   = "1144906123:AAGAwijc6ndXYO0R3bH9Rm0GHVyt0JOR5cM"
	WebhookURL = "https://hse-go-2020-1.herokuapp.com/" // телеграм обязует использовать https
	// WebhookURL = "https://83592a23.ngrok.io" // сслылка из телеграмма на локалхост
)

// источник из которого мы берем информацию
var rss = map[string]string{
	"Habr": "https://habrahabr.ru/rss/best/",
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func getNews(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	rss := new(RSS)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken) // инициализация BotAPI
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true // отладка
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL)) // обращение всех уведомлений на WebhookURL
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is working"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :8080")

	// получаем все обновления из канала updates
	for update := range updates {
		log.Printf("upd: %#v\n", update)
		url, ok := rss[update.Message.Text]
		if !ok {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				`there is only Habr feed availible`,
			))

			msg.ReplyMarkup = &tgbotapi.ReplyKeyboardMarkup{
				Keyboard: [][]tgbotapi.KeyboardButton{
					{
						Text: "Habr",
					},
				},
			}
			bot.Send(msg)
			continue
		}

		rss, err := getNews(url)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"sorry, error happend",
			))
		}
		for _, item := range rss.Items {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				item.URL+"\n"+item.Title, // item.URL - рабочая ссылка для перехода
			))
		}
	}
}
