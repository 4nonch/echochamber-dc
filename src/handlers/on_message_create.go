package handlers

import (
	"log"

	"github.com/4nonch/echochamber-dc/src/actions"
	"github.com/4nonch/echochamber-dc/src/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	_notMemberMsg = utils.MakeLocaleMap(
		"You're not a member of the channel.",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Вы не являетесь участником группы.",
		},
	)
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

	if !actions.CouldViewChannel(s, m) {
		msg := utils.GetLocalized(_notMemberMsg, discordgo.Locale(m.Author.Locale))
		actions.SendMessage(msg, s, m)
		return
	}

	actions.RedirectMessage(s, m)
}
