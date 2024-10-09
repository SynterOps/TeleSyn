package main

import (
	"fmt"
	"log"
	"os"

	"github.com/synternet/data-layer-sdk/pkg/dotenv"
)

func main() {
	// Load environment variables from .env file
	if err := dotenv.Load(); err != nil {
		log.Fatal("failed to load .env file:", err)
	}

	// Get wallet address from environment variable
	walletAddress := os.Getenv("WALLET_ADDRESS")
	if walletAddress == "" {
		log.Fatal("WALLET_ADDRESS environment variable not set")
	}

	// Initialize Telegram bot
	bot, err := telegram.NewBot()
	if err != nil {
		log.Fatal("failed to initialize Telegram bot:", err)
	}

	// Create a new command
	command := telegram.NewCommand{
		Name:        "wallet",
		Description: "Get your wallet address",
	}

	// Add a handler for the command
	command.Handler = func(update telegram.Update, args []string) {
		// Send a message with the wallet address
		bot.SendMessage(update.ChatID, fmt.Sprintf("Your wallet address is: %s", walletAddress), nil)
	}

	// Add the command to the bot
	bot.AddCommand(command)

	// Start the bot
	bot.Start()
}
