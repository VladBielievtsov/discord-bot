package main

import (
	"fmt"
	"go-discord-bot/internal/commands"
	"go-discord-bot/internal/config"
	"go-discord-bot/internal/types"
	"go-discord-bot/internal/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var songQueue = &types.Queue{}

func main() {
	err := godotenv.Load()
	utils.ErrorHandler(err)

	config.LoadConfig()

	sess, err := discordgo.New("Bot " + config.BotToken)
	utils.ErrorHandler(err)

	err = sess.Open()
	utils.ErrorHandler(err)
	defer sess.Close()

	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "hello",
			Description: "Say hello.",
		},
		{
			Name:        "hangman",
			Description: "Play Hangman",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "language",
					Description: "Select a language.",
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
		{
			Name:        "calculate",
			Description: "Calculate a math expression.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "expression",
					Description: "The math expression to calculate",
					Required:    true,
				},
			},
		},
		{
			Name:        "play",
			Description: "🎧 Add a song to the queue from a link or title.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "song",
					Description: "Link or title of the song.",
					Required:    true,
				},
			},
		},
	}

	for _, v := range cmds {
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, config.GuildID, v)
		utils.ErrorHandler(err)
	}

	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionMessageComponent:
			commands.HandleButtonInteraction(s, i)
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "hello":
				commands.Hello(s, i)
			case "hangman":
				commands.HangmanGame(s, i)
			case "calculate":
				commands.Calc(s, i)
			case "play":
				commands.Play(s, i, songQueue)
			}
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	fmt.Println("The bot is Online! Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
