package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	_ "github.com/synternet/data-layer-sdk/pkg/user"
)

func handleUserCommands(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Text == "/user" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User command received")
		bot.Send(msg)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("Message received from %s: %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I am the Synternet Trade Bot. I can help you buy and sell goods and services.")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello and welcome to the Synternet Trade Bot")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}

			handleUserCommands(bot, update)
		}
	}
}
