package handlers

import (
	"encoding/json"
	"go-discord-bot/db"
	"go-discord-bot/utils"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Answers struct {
	OriginChannelID string
	FavFood         string
	FavGame         string
	RecordID        int64
}

func (a *Answers) ToMessageEmbed() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Favorite food",
			Value: a.FavFood,
		},
		{
			Name:  "Favorite game",
			Value: a.FavGame,
		},
		{
			Name:  "Record ID",
			Value: strconv.FormatInt(a.RecordID, 10),
		},
	}

	return discordgo.MessageEmbed{
		Title:  "New responses!",
		Fields: fields,
	}
}

var Responses map[string]Answers = map[string]Answers{}

func UserPromptResponseHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	answers, ok := Responses[m.ChannelID]
	if !ok {
		return
	}

	if answers.FavFood == "" {
		answers.FavFood = m.Content

		s.ChannelMessageSend(m.ChannelID, "Great! What's your favorite game now? 🎮")
		Responses[m.ChannelID] = answers
		return
	} else {
		answers.FavGame = m.Content

		jbytes, err := json.Marshal(answers)
		utils.ErrorHandler(err)
		lastInserted, err := db.AddDiscordMessages(jbytes, m.ChannelID)
		utils.ErrorHandler(err)
		answers.RecordID = lastInserted

		s.ChannelMessageSend(m.ChannelID, "Great! Thanks you! 😊")
		embed := answers.ToMessageEmbed()
		s.ChannelMessageSendEmbed(answers.OriginChannelID, &embed)

		delete(Responses, m.ChannelID)
	}
}

func AnswersHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	spl := strings.Split(m.Content, " ")

	if len(spl) < 3 {
		s.ChannelMessageSend(m.ChannelID, "an ID must be provided. Ex: `!gobot answers 1`")
		return
	}

	id, err := strconv.Atoi(spl[2])
	utils.ErrorHandler(err)

	var recordID int64
	var answersStr string
	var userID int64

	query := "SELECT * FROM discord_messages WHERE id = $1"
	row := db.DB.QueryRow(query, id)
	err = row.Scan(&recordID, &answersStr, &userID)
	utils.ErrorHandler(err)

	var answers Answers
	err = json.Unmarshal([]byte(answersStr), &answers)
	utils.ErrorHandler(err)

	answers.RecordID = recordID
	embed := answers.ToMessageEmbed()
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func HelloWorldHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Hello World!")
}

func ProverbsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	proverbs := []string{
		"Don't communicate by sharing memory, share memory by communicating.",
		"Concurrency is not parallelism.",
		"Channels orchestrate; mutexes serialize.",
		"The bigger the interface, the weaker the abstraction.",
		"Make the zero value useful.",
		"interface{} says nothing.",
		"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
		"A little copying is better than a little dependency.",
		"Syscall must always be guarded with build tags.",
		"Cgo must always be guarded with build tags.",
		"Cgo is not Go.",
		"With the unsafe package there are no guarantees.",
		"Clear is better than clever.",
		"Reflection is never clear.",
		"Errors are values.",
		"Don't just check errors, handle them gracefully.",
		"Design the architecture, name the components, document the details.",
		"Documentation is for users.",
		"Don't panic.",
	}

	selection := rand.Intn(len(proverbs))

	author := discordgo.MessageEmbedAuthor{
		Name: "Rob Pike",
		URL:  "https://go-proverbs.github.io",
	}
	embed := discordgo.MessageEmbed{
		Title:  proverbs[selection],
		Author: &author,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func UserPromptHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	utils.ErrorHandler(err)

	if _, ok := Responses[channel.ID]; !ok {
		Responses[channel.ID] = Answers{
			OriginChannelID: m.ChannelID,
			FavFood:         "",
			FavGame:         "",
		}
		s.ChannelMessageSend(channel.ID, "Hey there! Here are some questions")
		s.ChannelMessageSend(channel.ID, "What's your favorite food?")
	} else {
		s.ChannelMessageSend(channel.ID, "We're still waiting... 😅")
	}
}
