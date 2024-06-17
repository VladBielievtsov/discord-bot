package main

import (
	"fmt"
	"go-discord-bot/handlers"
	"go-discord-bot/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	utils.ErrorHandler(err)

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		fmt.Println("BOT_TOKEN is not set in the environment variables")
		return
	}

	sess, err := discordgo.New("Bot " + token)
	utils.ErrorHandler(err)

	err = sess.Open()
	utils.ErrorHandler(err)
	defer sess.Close()

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "hello",
			Description: "Say hello",
		},
		{
			Name:        "hangman",
			Description: "Play Hangman",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "language",
					Description: "Select a language",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "English",
							Value: "en",
						},
						{
							Name:  "Ukraine",
							Value: "ua",
						},
					},
				},
			},
		},
	}

	guildID := os.Getenv("Guild_ID")

	for _, v := range commands {
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, guildID, v)
		utils.ErrorHandler(err)
	}

	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Name {
		case "hello":
			handlers.HelloHandler(s, i)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	fmt.Println("The bot is Online! Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
