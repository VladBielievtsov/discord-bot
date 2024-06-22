package handlers

import (
	"fmt"
	"go-discord-bot/utils"
	"time"

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

func PlayHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	songLink := i.ApplicationCommandData().Options[0].StringValue()

	guildID := i.GuildID

	guild, err := s.State.Guild(guildID)
	if err != nil {
		respondWithError(s, i, "Error fetching guild: "+err.Error())
		return
	}

	var userVoiceState *discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			userVoiceState = vs
			break
		}
	}

	if userVoiceState == nil {
		respondWithError(s, i, "You must be in a voice channel to use this command")
		return
	}

	vc, err := s.ChannelVoiceJoin(guildID, userVoiceState.ChannelID, false, true)
	if err != nil {
		respondWithError(s, i, "Failed to join voice channel: "+err.Error())
		return
	}
	defer vc.Disconnect()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Preparing to play: " + songLink,
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
		return
	}

	// playAudio()
	time.Sleep(5 * time.Second)

	s.ChannelMessageSend(i.ChannelID, "Finished playing")
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, errorMsg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error: " + errorMsg,
		},
	})
}
