package commands

import (
	"fmt"
	"go-discord-bot/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func Hello(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hello %v", i.Member.User.Username),
		},
	})
	utils.ErrorHandler(err)
}
