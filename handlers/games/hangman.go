package games

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var desc = [6]string{
	"```|‾‾‾‾‾‾|\n|\n|\n|\n|\n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|\n|\n|\n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      🙄 \n|\n|\n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      😟 \n|      👕 \n|\n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      😧 \n|      👕 \n|      🩳 \n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      😵 \n|      👕 \n|      🩳 \n|     👞👞 \n|\n|__________                      ```",
}

var (
	buttons = []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "A",
			Style:    discordgo.PrimaryButton,
			CustomID: "A",
		},
		discordgo.Button{
			Label:    "B",
			Style:    discordgo.PrimaryButton,
			CustomID: "B",
		},
		discordgo.Button{
			Label:    "C",
			Style:    discordgo.PrimaryButton,
			CustomID: "C",
		},
		discordgo.Button{
			Label:    "D",
			Style:    discordgo.PrimaryButton,
			CustomID: "D",
		},
		discordgo.Button{
			Label:    "Stop",
			Style:    discordgo.DangerButton,
			CustomID: "Stop",
		},
	}
	guessedLetters = []string{}
	steps          = 0
	word           string
)

func HangmanGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedLang := i.ApplicationCommandData().Options[0].StringValue()

	word = Words[rand.Intn(len(Words))]

	blanks := []string{}
	for range word {
		blanks = append(blanks, "🔵")
	}

	fmt.Println(word)

	var embed discordgo.MessageEmbed

	switch selectedLang {
	case "en":
		embed = discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    i.Member.User.Username,
				IconURL: i.Member.User.AvatarURL(""),
			},
			Title:       "Hangman",
			Color:       0xFFA500,
			Description: desc[steps],
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  fmt.Sprintf("Word (%d)", len(word)),
					Value: strings.Join(blanks, ""),
				},
			},
		}
	case "ua":
		embed = discordgo.MessageEmbed{
			Title: "Гра у Вісіліцю",
			Color: 0xFFA500,
		}
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unsupported language!",
			},
		})
		return
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: buttons,
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func HandleButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID

		switch customID {
		case "A", "B", "C", "D":
			guessedLetters = append(guessedLetters, customID)

			steps = steps + 1

			if steps == len(desc)-1 {
				Lost(s, i)
				return
			}

			embed := discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name:    i.Member.User.Username,
					IconURL: i.Member.User.AvatarURL(""),
				},
				Title:       "Hangman",
				Color:       0xFFA500,
				Description: desc[steps],
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Letters Guessed",
						Value: fmt.Sprintf("`%s`", strings.Join(guessedLetters, ", ")),
					},
				},
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{&embed},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: buttons,
						},
					},
				},
			})
			if err != nil {
				log.Printf("Error updating message: %v", err)
			}
		case "Stop":
			Lost(s, i)
		default:
			fmt.Printf("Unhandled button click with ID %s\n", customID)
		}
	}
}

func Lost(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.User.Username,
			IconURL: i.Member.User.AvatarURL(""),
		},
		Title:       "Hangman",
		Color:       0xFFA500,
		Description: desc[steps],
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Letters Guessed",
				Value: fmt.Sprintf("`%s`", strings.Join(guessedLetters, ", ")),
			},
			{
				Name:  "Game Over",
				Value: fmt.Sprintf("You lost! The word was **%s**.", word),
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
	if err != nil {
		log.Printf("Error updating message: %v", err)
	}
}

// {
// 	Name:  "Letters Guessed",
// 	Value: "",
// },
// {
// 	Name:  "Word (9)",
// 	Value: "🔵 🔵 🇪 🔵 🔵 🔵 🔵 🔵 🔵",
// },
// {
// 	Name:  "Game Over",
// 	Value: "You lost! The word was **partner**.",
// },
// {
// 	Name:  "Game Over",
// 	Value: "You won! The word was **partner**.",
// }
