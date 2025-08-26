package actions

import (
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

func RedirectMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(vars.ChannelID, m.Content)
}
