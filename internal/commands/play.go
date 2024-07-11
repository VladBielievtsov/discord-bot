package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
)

func Play(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	DownloadAudio(s, i, songLink)

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

func DownloadAudio(s *discordgo.Session, i *discordgo.InteractionCreate, videoURL string) {
	// Download Audio
	client := youtube.Client{}

	video, err := client.GetVideo(videoURL)
	if err != nil {
		respondWithError(s, i, "Error getting video: "+err.Error())
		return
	}

	var audioFormat *youtube.Format
	for _, format := range video.Formats {
		if format.ItagNo == 140 {
			audioFormat = &format
			break
		}
	}

	if audioFormat == nil {
		respondWithError(s, i, "Audio format not found")
		return
	}

	tempDir := "temp"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		respondWithError(s, i, "Error creating directory: "+err.Error())
		return
	}

	id := uuid.New().String()
	fileName := id + ".m4a"

	outFilePath := filepath.Join(tempDir, fileName)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		respondWithError(s, i, "Error creating file: "+err.Error())
		return
	}
	defer outFile.Close()

	stream, _, err := client.GetStream(video, audioFormat)
	if err != nil {
		respondWithError(s, i, "Error getting stream: "+err.Error())
		return
	}
	defer stream.Close()

	_, err = io.Copy(outFile, stream)
	if err != nil {
		respondWithError(s, i, "Error saving audio: "+err.Error())
		return
	}
}
