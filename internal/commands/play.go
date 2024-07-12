package commands

import (
	"fmt"
	"go-discord-bot/internal/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
)

func Play(s *discordgo.Session, i *discordgo.InteractionCreate, queue *types.Queue) {
	songLink := i.ApplicationCommandData().Options[0].StringValue()

	queueItem := types.QueueItem{
		ID:       uuid.New().String(),
		VideoURL: songLink,
		User: types.QueueUser{
			Name:      i.Member.User.Username,
			AvatarURL: i.Member.User.AvatarURL(""),
		},
	}

	queue.Enqueue(queueItem)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Preparing to play: " + songLink,
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
		return
	}

	if len(queue.Items) == 1 {
		go PlayNextInQueue(s, i, queue)
	}

}

var vc *discordgo.VoiceConnection

func PlayNextInQueue(s *discordgo.Session, i *discordgo.InteractionCreate, queue *types.Queue) {
	if queue.IsEmpty() {
		if vc != nil {
			vc.Disconnect()
			vc = nil
		}
		return
	}

	queueItem := queue.Dequeue()
	if queueItem == nil {
		return
	}

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

	if vc == nil {
		vc, err = s.ChannelVoiceJoin(guildID, userVoiceState.ChannelID, false, true)
		if err != nil {
			respondWithError(s, i, "Failed to join voice channel: "+err.Error())
			return
		}
	}

	audioFile, err := DownloadAudio(queueItem.VideoURL)
	if err != nil {
		respondWithError(s, i, "Error downloading audio: "+err.Error())
		return
	}
	fmt.Println("TEst 1")
	err = PlayAudioFile(vc, audioFile)
	if err != nil {
		respondWithError(s, i, "Error playing audio: "+err.Error())
		return
	}

	fmt.Println("TEst 2")

	// time.Sleep(5 * time.Second)

	s.ChannelMessageSend(i.ChannelID, "Finished playing")

	PlayNextInQueue(s, i, queue)
}

func DownloadAudio(videoURL string) (string, error) {
	// Download Audio
	client := youtube.Client{}

	video, err := client.GetVideo(videoURL)
	if err != nil {
		return "", fmt.Errorf("error getting video: %w", err)
	}

	var audioFormat *youtube.Format
	for _, format := range video.Formats {
		if format.ItagNo == 140 {
			audioFormat = &format
			break
		}
	}

	if audioFormat == nil {
		return "", fmt.Errorf("Audio format not found")
	}

	tempDir := "temp"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %w", err)
	}

	id := uuid.New().String()
	fileName := id + ".mp3"

	outFilePath := filepath.Join(tempDir, fileName)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer outFile.Close()

	stream, _, err := client.GetStream(video, audioFormat)
	if err != nil {
		return "", fmt.Errorf("error getting stream: %w", err)
	}
	defer stream.Close()

	_, err = io.Copy(outFile, stream)
	if err != nil {
		return "", fmt.Errorf("error saving audio: %w", err)
	}

	opusFilePath := filepath.Join(tempDir, id+".opus")
	cmd := exec.Command("ffmpeg", "-i", outFilePath, "-c:a", "libopus", "-b:a", "64k", opusFilePath)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error converting to DCA: %w", err)
	}

	return opusFilePath, nil
}

func PlayAudioFile(vc *discordgo.VoiceConnection, filename string) error {
	fmt.Println("Starting to play audio...")

	err := vc.Speaking(true)
	if err != nil {
		return fmt.Errorf("failed to turn on mic: %w", err)
	}
	defer func() {
		vc.Speaking(false)
		fmt.Println("Stopped playing audio.")
	}()

	dcaFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening DCA file: %w", err)
	}
	defer dcaFile.Close()

	decoder := dca.NewDecoder(dcaFile)
	if decoder == nil {
		return fmt.Errorf("failed to create DCA decoder")
	}

	done := make(chan error)

	go func() {
		defer close(done)
		for {
			frame, err := decoder.OpusFrame()
			if err == io.EOF {
				done <- nil
				return
			} else if err != nil {
				done <- fmt.Errorf("error decoding Opus frame: %w", err)
				return
			}

			if !vc.Ready || vc.OpusSend == nil {
				fmt.Println("Voice connection not ready or OpusSend channel is nil")
				continue
			}

			vc.OpusSend <- frame
			fmt.Println("Sent Opus frame")
		}
	}()

	for err := range done {
		if err != nil {
			fmt.Println("Playback error:", err)
			return fmt.Errorf("playback error: %w", err)
		}
	}

	fmt.Println("Playback finished successfully")

	return nil
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, errorMsg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error: " + errorMsg,
		},
	})
}
