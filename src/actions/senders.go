package actions

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func SendInteractionMessage(content string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		},
	)
	if err != nil {
		log.Printf("Failed to send interaction message \"%v\": %v", content, err)
	}
}

func SendMessage(content string, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, content)
	if err != nil {
		log.Printf("Failed to send message \"%v\": %v", content, err)
	}
}
