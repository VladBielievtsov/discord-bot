package handlers

import (
	"fmt"
	"go-discord-bot/utils"

	"github.com/Knetic/govaluate"
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

// FIXME: If passing a non-mathematical expression like "hello" causes an error

func CalcHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	expressionStr := i.ApplicationCommandData().Options[0].StringValue()

	expression, err := govaluate.NewEvaluableExpression(expressionStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s = %v", expression, result),
		},
	})
	utils.ErrorHandler(err)
}
