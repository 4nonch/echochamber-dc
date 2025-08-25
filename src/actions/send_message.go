package actions

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func SendMessage(content string, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, content)
	if err != nil {
		log.Printf("Failed to send message \"%v\": %v", content, err)
	}
}
