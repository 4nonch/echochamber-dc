package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(c *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf(
		"New message: (type=%d) (attachments=%d) %s",
		m.Type,
		len(m.Attachments),
		m.Content,
	)
}
