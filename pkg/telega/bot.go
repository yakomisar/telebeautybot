package telega

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var baseMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("✏️ Записаться"),
		tgbotapi.NewKeyboardButton("ℹ️ О нас"),
	),
)

var manicureService = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Педикюр", "Педикюр"),
		tgbotapi.NewInlineKeyboardButtonData("Маникюр", "Маникюр"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Закончить выбор", "Закончить выбор"),
	),
)

const (
	commandStart  = "start"
	replyStart    = `Вас приветствует бот салона красоты "BonjourManicure"`
	commandInfo   = "info"
	replyInfo     = "ВРЕМЯ РАБОТЫ\nПн.-Вс.\n10.00 - 21.00\n\n+7 (926) 227 77 50\n\nАДРЕС САЛОНА\nМосква, Знаменские садки д5/1 (МОДЖО)\n\nСледите за нами в телеграм:\n@bonjourmanicure_butovo"
	chooseService = "Выберите услугу на которую хотите записаться"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

func (b *Bot) handleMessages(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = baseMenu
	if msg.Text == "ℹ️ О нас" {
		loc := tgbotapi.NewLocation(message.Chat.ID, 55.571815, 37.570249)
		b.bot.Send(loc)
		msg.Text = replyInfo
		b.bot.Send(msg)
	} else if msg.Text == "✏️ Записаться" {
		msg.Text = chooseService
		b.bot.Send(msg)
		service := tgbotapi.NewMessage(message.Chat.ID, message.Text)
		service.ReplyToMessageID = message.MessageID
		service.ReplyMarkup = manicureService
	} else {
		msg.Text = replyStart
		b.bot.Send(msg)
	}
}

func (b *Bot) handleCommands(command *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(command.Chat.ID, command.Text)
	msg.ReplyToMessageID = command.MessageID
	msg.ReplyMarkup = baseMenu
	switch command.Command() {
	case commandStart:
		msg.Text = replyStart
		b.bot.Send(msg)
	case commandInfo:
		msg.ReplyToMessageID = command.MessageID
		loc := tgbotapi.NewLocation(command.Chat.ID, 55.571815, 37.570249)
		b.bot.Send(loc)
		msg.Text = replyInfo
		b.bot.Send(msg)
	default:
		msg.Text = "Unknown command"
		b.bot.Send(msg)
	}
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			b.handleCommands(update.Message)
			continue
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := b.bot.Send(msg); err != nil {
				panic(err)
			}
			continue
		}
		b.handleMessages(update.Message)
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	b.handleUpdates(updates)
	return nil
}
