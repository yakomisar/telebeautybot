package main

import (
	"log"
	"telebeautybot/pkg/telega"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5121500762:AAH1GLgug_B4ddj-Jb4oB_jR0xI1SMLU8L8")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	teleBot := telega.NewBot(bot)
	if err := teleBot.Start(); err != nil {
		log.Fatal(err)
	}
}
