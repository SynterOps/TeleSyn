package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/synternet/data-layer-sdk/pkg/dotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/service"
)

func main() {
	urls := flag.String("urls", os.Getenv("NATS_URL"), "NATS urls")
	source := flag.String("source", "synternet.price.single.ATOM", "Source subject to stream from")
	creds := flag.String("nats-creds", os.Getenv("NATS_CREDS"), "NATS credentials file")
	nkey := flag.String("nats-nkey", os.Getenv("NATS_NKEY"), "NATS NKey seed")
	jwt := flag.String("nats-jwt", os.Getenv("NATS_JWT"), "NATS JWT string")
	verbose := flag.Bool("verbose", false, "Verbose logs")

	flag.Parse()

	conn, err := options.MakeNatsConnection("Streaming consumer", *urls, *creds, *nkey, *jwt, "", "", "")
	if err != nil {
		panic(fmt.Errorf("failed to create NATS connection: %w", err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	opts := []options.Option{
		service.WithContext(ctx),
		service.WithName("StreamingConsumer"),
		service.WithNats(conn),
		service.WithVerbose(*verbose),
		service.WithParam(service.WithSource, *source),
		service.WithParam(service.WithNKey, *nkey),
	}

	consumer, err := service.NewStreamingConsumer(opts...)
	if err != nil {
		panic(fmt.Errorf("failed to create the consumer: %w", err))
	}

	consumer.Start()
	defer consumer.Close()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panicf("failed to initialize bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	updatesChannel, err := bot.GetUpdatesChannel(updates)
	if err != nil {
		log.Panicf("failed to get updates channel: %v", err)
	}

	for update := range updatesChannel {
		if update.Message != nil {
			log.Printf("Message received from %s: %s", update.Message.From.UserName, update.Message.Text)

			switch strings.ToLower(update.Message.Text) {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I am the Synternet Trade Bot. I can help you buy and sell goods and services. Use /start and /help for more info.")
				bot.Send(msg)
			case "/account":
				account, err := queryAccount(update.Message.Chat.ID, update.Message.MessageID, bot)
				if err != nil {
					log.Printf("Error querying account: %v", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There was an error retrieving your account details.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your account details:\nLamports: "+strconv.FormatInt(account.Lamports, 10))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "/price":
				price, err := queryPrice(consumer)
				if err != nil {
					logMultiplier:
					log.Printf("Error querying price: %v", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There was an error retrieving the price.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The current price is: "+price)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "/help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Available commands:\n/start - Start the bot\n/account - Get your account details\n/price - Get the current price\n/help - Display this help message")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I didn't understand that command. Use /help to see a list of available commands.")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}

type Account struct {
	Lamports int64
}

func queryAccount(chatID int64, messageID int, bot *tgbotapi.BotAPI) (account Account, err error) {
	// Placeholder implementation, replace with actual account query logic
	time.Sleep(2 * time.Second)
	return Account{
		Lamports: 1000000,
	}, nil
}

func queryPrice(consumer *service.StreamingConsumer) (price string, err error) {
	// Placeholder implementation, replace with actual price query logic
	msg, err := consumer.Request("synternet.price.single.ATOM", nil, 2*time.Second)
	if err != nil {
		return "", err
	}
	return string(msg.Data), nil
}
