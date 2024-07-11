package commands

import (
	"fmt"
	"go-discord-bot/internal/utils"

	"github.com/Knetic/govaluate"
	"github.com/bwmarrin/discordgo"
)

// FIXME: If passing a non-mathematical expression like "hello" causes an error

func Calc(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
