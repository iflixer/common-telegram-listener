package listener

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"telegram-listener/database"
	"telegram-listener/reaction"
	"time"

	"gopkg.in/telebot.v4"
)

type FlixBot struct {
	telegramBot database.TelegramBot
	TgBot       *telebot.Bot
}

func (flixBot *FlixBot) Register(dbService *database.Service, reactionService *reaction.Service) (err error) {
	log.Println("Registering bot:", flixBot.telegramBot.ID, flixBot.telegramBot.Name)
	pref := telebot.Settings{
		Token:  flixBot.telegramBot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	log.Println("Bot listen URL:", flixBot.telegramBot.ListenURL)
	log.Println("Webhook info:", flixBot.getWebhookInfo(flixBot.telegramBot.Token))

	flixBot.TgBot, err = telebot.NewBot(pref)
	if err != nil {
		if strings.Contains(err.Error(), "Conflict") {
			log.Fatalf("❌ Уже другой бот слушает webhook: %v", err)
		}
		log.Println("Failed to register bot:", flixBot.telegramBot.ID, flixBot.telegramBot.Name, err)
		return
	}

	log.Println("Starting bot:", flixBot.telegramBot.ID, flixBot.telegramBot.Name)

	reactionService.RegisterReactions(flixBot.telegramBot.ID, flixBot.telegramBot.AppURL, flixBot.TgBot)

	// flixBot.TgBot.Handle("/start", func(c telebot.Context) error {
	// 	return c.Send("Hello, I am your bot!")
	// })
	go flixBot.TgBot.Start()

	return
}

func (c *FlixBot) Stop() {
	log.Println("Stopping bot:", c.telegramBot.ID, c.telegramBot.Name)
	if c.TgBot != nil {
		c.TgBot.Stop()
	}
}

func (c *FlixBot) getWebhookInfo(botToken string) string {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getWebhookInfo", botToken)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Неверный код ответа: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	// Печатаем JSON-ответ Telegram
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult)
}
