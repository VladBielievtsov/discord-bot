package main

import (
	"fmt"
	"go-discord-bot/db"
	"go-discord-bot/handlers"
	"go-discord-bot/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const prefix string = "!gobot"

func main() {
	err := godotenv.Load()
	utils.ErrorHandler(err)

	token := os.Getenv("BOT_TOKEN")
	sess, err := discordgo.New("Bot " + token)
	utils.ErrorHandler(err)

	db.ConnectDatabase()
	defer db.DB.Close()

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		// DM Logic
		if m.GuildID == "" {
			handlers.UserPromptResponseHandler(s, m)
		}

		// Server Logic

		args := strings.Split(m.Content, " ")

		if args[0] != prefix {
			return
		}

		if args[1] == "Hello" {
			handlers.HelloWorldHandler(s, m)
		}

		if args[1] == "proverbs" {
			handlers.ProverbsHandler(s, m)
		}

		if args[1] == "prompt" {
			handlers.UserPromptHandler(s, m)
		}

		if args[1] == "answers" {
			handlers.AnswersHandler(s, m)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	utils.ErrorHandler(err)
	defer sess.Close()

	fmt.Println("The bot is Online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
