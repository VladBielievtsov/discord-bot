package main

import (
	"fmt"
	"go-discord-bot/db"
	"go-discord-bot/handlers"
	"go-discord-bot/utils"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var images = [8]string{
	"https://i.pinimg.com/736x/e0/e2/75/e0e2751ec0a039f536234e3f87432acf.jpg",
	"https://i.pinimg.com/564x/b0/ff/44/b0ff44d634233716697b774da9a4ad7a.jpg",
	"https://i.pinimg.com/736x/fc/9b/c0/fc9bc08b6313c642965c491b94010e75.jpg",
	"https://i.pinimg.com/736x/4e/56/d8/4e56d8cebeceeab0fc5fe2fad198d6ba.jpg",
	"https://i.pinimg.com/564x/bf/6f/8a/bf6f8ab7b14f71af53aa226c40efdb86.jpg",
	"https://i.pinimg.com/564x/4f/c2/45/4fc24581a7e7349e4551ad1924918d15.jpg",
	"https://i.pinimg.com/564x/bd/26/b3/bd26b3b9772adb0ac0854cef3ebd16bf.jpg",
	"https://i.pinimg.com/564x/e9/f8/dc/e9f8dc62e3ad3053ca9aa195edd72016.jpg",
}

const prefix string = ".cat"

func main() {
	err := godotenv.Load()
	utils.ErrorHandler(err)

	token := os.Getenv("BOT_TOKEN")
	sess, err := discordgo.New("Bot " + token)
	utils.ErrorHandler(err)

	db.ConnectDatabase()
	defer db.DB.Close()

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.Emoji.Name == "🔥" {
			s.GuildMemberRoleAdd(r.GuildID, r.UserID, "1251984426032824400")
			s.ChannelMessageSend(r.ChannelID, fmt.Sprintf("%v has been added to %v", r.UserID, r.Emoji.Name))
		}
	})

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
		if r.Emoji.Name == "🔥" {
			s.GuildMemberRoleRemove(r.GuildID, r.UserID, "1251984426032824400")
			s.ChannelMessageSend(r.ChannelID, fmt.Sprintf("%v has been removed from %v", r.UserID, r.Emoji.Name))
		}
	})

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

		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Enter `.cat help` to see all commands")
			return
		}

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

		if args[1] == "rand" {
			randImage := rand.Intn(len(images))
			s.ChannelMessageSend(m.ChannelID, images[randImage])
		}

		if args[1] == "help" {
			handlers.HelpHandler(s, m)
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
