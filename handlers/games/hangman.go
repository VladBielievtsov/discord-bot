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
	letters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	buttons    = generateButtons(letters)
	stopButton = discordgo.Button{
		Label:    "Stop",
		Style:    discordgo.DangerButton,
		CustomID: "Stop",
	}
	nextButton = discordgo.Button{
		Style:    discordgo.SuccessButton,
		CustomID: "Next",
		Emoji: &discordgo.ComponentEmoji{
			Name: "➡️",
		},
	}
	prevButton = discordgo.Button{
		Style:    discordgo.SuccessButton,
		CustomID: "Prev",
		Emoji: &discordgo.ComponentEmoji{
			Name: "⬅️",
		},
	}
	currentState   = "firstButtons"
	guessedLetters = []string{}
	steps          = 0
	word           string
	blanks         = []string{}
)

func HangmanGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedLang := i.ApplicationCommandData().Options[0].StringValue()

	word = Words[rand.Intn(len(Words))]

	for range word {
		blanks = append(blanks, "🔵 ")
	}

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
			Title:       "Гра у Вісіліцю",
			Color:       0xFFA500,
			Description: "Українська версія ще не готова",
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
			Embeds:     []*discordgo.MessageEmbed{&embed},
			Components: generateFirstButtons(),
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func generateButtons(letters string) []discordgo.MessageComponent {
	var buttons []discordgo.MessageComponent
	for _, letter := range letters {
		buttons = append(buttons, discordgo.Button{
			Label:    string(letter),
			Style:    discordgo.PrimaryButton,
			CustomID: string(letter),
		})
	}
	return buttons
}

func generateFirstButtons() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: buttons[:4],
		},
		discordgo.ActionsRow{
			Components: buttons[4:8],
		},
		discordgo.ActionsRow{
			Components: buttons[8:12],
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				nextButton, stopButton,
			},
		},
	}
}

func generateSecondButtons() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: buttons[12:16],
		},
		discordgo.ActionsRow{
			Components: buttons[16:20],
		},
		discordgo.ActionsRow{
			Components: buttons[20:24],
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				prevButton, stopButton, buttons[24], buttons[25],
			},
		},
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
							Components: generateFirstButtons(),
						},
					},
				},
			})
			if err != nil {
				log.Printf("Error updating message: %v", err)
			}
		case "Stop":
			Lost(s, i)
		case "Next":
			switch currentState {
			case "firstButtons":
				currentState = "secondButtons"

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
							Name:  fmt.Sprintf("Word (%d)", len(word)),
							Value: strings.Join(blanks, ""),
						},
					},
				}

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds:     []*discordgo.MessageEmbed{&embed},
						Components: generateSecondButtons(),
					},
				})
				if err != nil {
					log.Printf("Error updating message: %v", err)
				}
			case "secondButtons":
				fmt.Println("Already in the second button state, handle accordingly")
			}
		case "Prev":
			switch currentState {
			case "secondButtons":
				currentState = "firstButtons"

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
							Name:  fmt.Sprintf("Word (%d)", len(word)),
							Value: strings.Join(blanks, ""),
						},
					},
				}

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds:     []*discordgo.MessageEmbed{&embed},
						Components: generateFirstButtons(),
					},
				})
				if err != nil {
					log.Printf("Error updating message: %v", err)
				}
			case "firstButtons":
				fmt.Println("Already in the second button state, handle accordingly")
			}
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
