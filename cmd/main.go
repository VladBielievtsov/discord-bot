package main

import (
	"fmt"
	"go-discord-bot/handlers"
	"go-discord-bot/handlers/games"
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

	guildID := os.Getenv("Guild_ID")
	if guildID == "" {
		fmt.Println("Guild_ID is not set in the environment variables")
		return
	}

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
						// {
						// 	Name:  "Ukraine",
						// 	Value: "ua",
						// },
					},
				},
			},
		},
	}

	for _, v := range commands {
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, guildID, v)
		utils.ErrorHandler(err)
	}

	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionMessageComponent:
			games.HandleButtonInteraction(s, i)
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "hello":
				handlers.HelloHandler(s, i)
			case "hangman":
				games.HangmanGame(s, i)
			}
		}
	})

	// sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	if i.Type == discordgo.InteractionMessageComponent {
	// 		handlers.HandleButtonInteraction(s, i)
	// 	}
	// })

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	fmt.Println("The bot is Online! Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
