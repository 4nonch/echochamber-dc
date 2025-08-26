package actions

import (
	"fmt"
	"log"

	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

func CouldViewChannel(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	guild, err := s.State.Guild(vars.GuildID)
	if err == nil && guild.OwnerID == m.Author.ID {
		return true
	}

	perms, err := GetChannelPermissions(s, m.Author.ID)
	if err != nil {
		msg := fmt.Sprintf("Unable to get user's permissions for channel \"%s\": %v", vars.ChannelID, err)
		log.Printf(msg)
		SendMessage(msg, s, m)
		return false
	}

	return (perms&discordgo.PermissionViewChannel != 0) ||
		(perms&discordgo.PermissionAdministrator != 0)
}
