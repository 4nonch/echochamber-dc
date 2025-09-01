package services

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/4nonch/echochamber-dc/src/actions"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

func RedirectMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !CouldViewChannel(s, m) {
		msg := "You're not a member of the channel."
		actions.SendMessage(msg, s, m)
		return
	}
	if err := ValidateContent(s, m); err != nil {
		actions.SendMessage(err.Error(), s, m)
	}
	if err := ValidateMedia(s, m); err != nil {
		actions.SendMessage(err.Error(), s, m)
	}

	stream, failed := fetchAttachments(s, m)
	if failed {
		return
	}

	ms, failed := prepareMessage(s, m, stream.Files)
	if failed {
		return
	}

	_, err := s.ChannelMessageSendComplex(vars.ChannelID, ms)
	stream.Close()

	if err != nil {
		msg := fmt.Sprintf(
			"Failed to redirect message (attachments: %d): %v",
			len(m.Attachments),
			err,
		)
		log.Println(msg)
		actions.SendMessage(msg, s, m)
		return
	}

	actions.SendMessage("Message successfully delivered.", s, m)
}

// Returns attachment files (if there is no one, Files attribute will be simply an empty slice).
// If second error is true - then error occurred and notified.
func fetchAttachments(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) (stream actions.StreamFiles, failed bool) {
	if len(m.Attachments) == 0 {
		return stream, false
	}

	var errs chan error
	stream, errs = actions.GetAttachments(m.Attachments)
	if len(errs) != 0 {
		for err := range errs {
			log.Println("Failed to get attachment: ", err)
		}
		actions.SendMessage("An error occurred while trying to download attachments.", s, m)
		return stream, true
	}
	return stream, false
}

// If failed (error occurred) - handles all errors and returns "true" boolean
func prepareMessage(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	fs []*discordgo.File,
) (ms *discordgo.MessageSend, failed bool) {
	reference, content, err := extractReference(m)
	if err != nil {
		msg := fmt.Sprintf("Failed to prepare message: %v", err)
		log.Println(msg)
		actions.SendMessage(msg, s, m)
		return nil, true
	}

	if len(m.Attachments) == 0 && strings.TrimSpace(content) == "" {
		if reference != nil {
			actions.SendMessage("Can't reply with an empty message.", s, m)
			return nil, true
		}
		actions.SendMessage("Can't send an empty message.", s, m)
		return nil, true
	}

	if utf8.RuneCountInString(m.Content) > vars.MaxMessageChars {
		actions.SendMessage("Message result is too big, failed to send.", s, m)
		return nil, true
	}

	return &discordgo.MessageSend{
		Content:   content,
		Files:     fs,
		Reference: reference,
	}, false
}
