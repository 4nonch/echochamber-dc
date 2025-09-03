package services

import (
	"errors"
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/4nonch/echochamber-dc/src/actions"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

// If you're owner - skip check.
// Be aware of setting misleading ChannelID, since it won't be checked if you're the owner.
func CouldViewChannel(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	guild, err := s.State.Guild(vars.GuildID)
	if err == nil && guild.OwnerID == m.Author.ID {
		return true
	}

	perms, err := actions.GetChannelPermissions(s, m.Author.ID)
	if err == nil {
		msg := fmt.Sprintf("Unable to get user's permissions for channel \"%s\": %v", vars.ChannelID, err)
		log.Printf(msg)
		return false
	}

	return (perms&discordgo.PermissionViewChannel != 0) ||
		(perms&discordgo.PermissionAdministrator != 0)
}

func ValidateContent(s *discordgo.Session, m *discordgo.MessageCreate) error {
	count := utf8.RuneCountInString(m.Content)
	if count > vars.MaxMessageChars {
		msg := fmt.Sprintf(
			"You're message is too big (%d). Maximum allowed size: %d characters.",
			count,
			vars.MaxMessageChars,
		)
		return errors.New(msg)
	}
	return nil
}

func ValidateMedia(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Max attachments count check
	if len(m.Attachments) > vars.MaxAttachmentsCount {
		msg := fmt.Sprintf(
			"Too much attachments (%d). Maximum allowed count: %d",
			len(m.Attachments),
			vars.MaxAttachmentsCount,
		)
		return errors.New(msg)
	}

	// Max attachments size check
	var total int
	for _, a := range m.Attachments {
		total += a.Size
	}
	if total > vars.MaxAttachmentsBytes {
		msg := fmt.Sprintf(
			"Attachments are too big (%.2f Mb). Maximum allowed size: %.2f Mb.",
			float32(total)/1024/1024,
			float32(vars.MaxAttachmentsBytes)/1024/1024,
		)
		return errors.New(msg)
	}
	return nil
}
