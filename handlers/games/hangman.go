package games

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/enescakir/emoji"
)

var desc = [6]string{
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      😵 \n|      👕 \n|      🩳 \n|     👞👞 \n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      😧 \n|      👕 \n|      🩳 \n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      😟 \n|      👕 \n|\n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|      🙄 \n|\n|\n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|      🎩 \n|\n|\n|\n|\n|\n|__________                      ```",
	"```|‾‾‾‾‾‾|\n|\n|\n|\n|\n|\n|\n|__________                      ```",
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
	currentState          = "firstButtons"
	guessedLetters        = []string{}
	lives                 = 5
	word                  string
	correctGuessedLetters = []string{}
	blanks                = []string{}
	words                 = []string{"golang", "js", "php"}
)

func HangmanGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedLang := i.ApplicationCommandData().Options[0].StringValue()

	lives = 5
	blanks = []string{}
	guessedLetters = []string{}
	correctGuessedLetters = []string{}
	currentState = "firstButtons"

	word = words[rand.Intn(len(words))]

	for range word {
		blanks = append(blanks, "🔵 ")
		correctGuessedLetters = append(correctGuessedLetters, "_")
	}

	var embed discordgo.MessageEmbed

	switch selectedLang {
	case "en":
		embed = generateEmbed(i, []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Word (%d)", len(word)),
				Value: strings.Join(blanks, ""),
			},
		})
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

func generateEmbed(i *discordgo.InteractionCreate, fields []*discordgo.MessageEmbedField) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.User.Username,
			IconURL: i.Member.User.AvatarURL(""),
		},
		Title:       "Hangman",
		Color:       0xFFA500,
		Description: desc[lives],
		Fields:      fields,
	}
}

func HandleButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID

		switch customID {
		case "Stop":
			Lost(s, i)
		case "Next":
			switch currentState {
			case "firstButtons":
				currentState = "secondButtons"

				embed := generateEmbed(i, []*discordgo.MessageEmbedField{
					{
						Name:  "Letters Guessed",
						Value: fmt.Sprintf("`%s`", strings.Join(guessedLetters, ", ")),
					},
					{
						Name:  fmt.Sprintf("Word (%d)", len(word)),
						Value: strings.Join(blanks, ""),
					},
				})

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

				embed := generateEmbed(i, []*discordgo.MessageEmbedField{
					{
						Name:  "Letters Guessed",
						Value: fmt.Sprintf("`%s`", strings.Join(guessedLetters, ", ")),
					},
					{
						Name:  fmt.Sprintf("Word (%d)", len(word)),
						Value: strings.Join(blanks, ""),
					},
				})

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
			if len(customID) == 1 && customID[0] >= 'A' && customID[0] <= 'Z' {
				guessedLetters = append(guessedLetters, customID)
				correctGuess := false
				for i, wordLetter := range word {
					if strings.ToLower(customID) == string(wordLetter) {
						blanks[i] = emoji.Parse(":regional_indicator_" + strings.ToLower(customID) + ":")
						correctGuessedLetters[i] = strings.ToLower(customID)
						correctGuess = true
					}
				}

				if !correctGuess {
					lives--
				}

				if lives == 0 {
					Lost(s, i)
					return
				}

				if strings.Join(correctGuessedLetters, "") == word {
					embed := generateEmbed(i, []*discordgo.MessageEmbedField{
						{
							Name:  "Letters Guessed",
							Value: fmt.Sprintf("`%s`", strings.Join(guessedLetters, ", ")),
						},
						{
							Name:  "Game Over",
							Value: fmt.Sprintf("You won! The word was **%s**.", word),
						},
					})

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseUpdateMessage,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{&embed},
						},
					})
					if err != nil {
						log.Printf("Error updating message: %v", err)
					}
					return
				}

				embed := generateEmbed(i, []*discordgo.MessageEmbedField{
					{
						Name:  "Letters Guessed",
						Value: fmt.Sprintf("`%s`", strings.Join(guessedLetters, ", ")),
					},
					{
						Name:  fmt.Sprintf("Word (%d)", len(word)),
						Value: strings.Join(blanks, ""),
					},
				})

				switch currentState {
				case "firstButtons":
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
				case "secondButtons":
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
				}

			} else {
				fmt.Printf("Unhandled button click with ID %s\n", customID)
			}
		}
	}
}

func Lost(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := generateEmbed(i, []*discordgo.MessageEmbedField{
		{
			Name:  "Letters Guessed",
			Value: fmt.Sprintf("`%s`", strings.Join(guessedLetters, ", ")),
		},
		{
			Name:  "Game Over",
			Value: fmt.Sprintf("You lost! The word was **%s**.", word),
		},
	})

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
