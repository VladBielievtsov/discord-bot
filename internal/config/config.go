package config

import (
	"fmt"
	"log"
	"os"
)

var (
	BotToken string
	GuildID  string
)

func LoadConfig() {
	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		log.Println("BOT_TOKEN environment variable not set")
		return
	}

	GuildID := os.Getenv("Guild_ID")
	if GuildID == "" {
		fmt.Println("Guild_ID is not set in the environment variables")
		return
	}
}
