package main

import (
	"log"
	"telebeautybot/pkg/telega"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	teleBot := telega.NewBot(bot)
	if err := teleBot.Start(); err != nil {
		log.Fatal(err)
	}
}
