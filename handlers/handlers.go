package handlers

import (
	"fmt"
	"go-discord-bot/utils"

	"github.com/bwmarrin/discordgo"
)

func HelloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hello %v", i.Member.User.Username),
		},
	})
	utils.ErrorHandler(err)
}
