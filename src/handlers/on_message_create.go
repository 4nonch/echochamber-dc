package handlers

import (
	"log"

	"github.com/4nonch/echochamber-dc/src/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.add(OnMessageCreate)
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID != "" ||
		m.Author.ID == s.State.User.ID {
		return
	}

	log.Printf(
		"New message: (type=%d) (attachments=%d) %s",
		m.Type,
		len(m.Attachments),
		m.Content,
	)

	services.RedirectMessage(s, m)
}
