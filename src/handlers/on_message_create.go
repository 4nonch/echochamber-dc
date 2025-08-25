package handlers

import (
	"log"

	"github.com/4nonch/echochamber-dc/src/actions"
	"github.com/bwmarrin/discordgo"
)

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
	actions.SendMessage("OK, IT WORKS", s, m)
}
