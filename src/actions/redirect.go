package actions

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/4nonch/echochamber-dc/src/patterns"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

// Errors
var (
	errEmptyMessage = errors.New("Can't send an empty message.")
)

func RedirectMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !HasValidContent(s, m) || !HasValidMedia(s, m) {
		return
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
			"Failed to redirect message \"%s\" (attachments: %d): %v",
			m.Content,
			len(m.Attachments),
			err,
		)
		log.Println(msg)
		SendMessage(msg, s, m)
		return
	}

	SendMessage("Message successfully delivered.", s, m)
}

// Returns attachment files (if there is no one, Files attribute will be simply an empty slice).
// If second error is true - then error occurred and notified.
func fetchAttachments(s *discordgo.Session, m *discordgo.MessageCreate) (stream StreamFiles, failed bool) {
	if len(m.Attachments) == 0 {
		return stream, false
	}

	var errs chan error
	stream, errs = GetAttachments(m.Attachments)
	if len(errs) != 0 {
		for err := range errs {
			log.Println("Failed to get attachment: ", err)
		}
		SendMessage("An error occurred while trying to download attachments.", s, m)
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
	reference, content, err := extractReference(m.Content)
	if errors.Is(err, errEmptyMessage) {
		SendMessage("Can't reply with an empty message.", s, m)
		return nil, true
	} else if err != nil {
		msg := fmt.Sprintf("Failed to prepare message: %v", err)
		log.Println(msg)
		SendMessage(msg, s, m)
		return nil, true
	}

	return &discordgo.MessageSend{
		Content:   content,
		Files:     fs,
		Reference: reference,
	}, false
}

// Return MessageReference if original content contained a reply.
// Link to the replied message will be removed from original content
func extractReference(c string) (*discordgo.MessageReference, string, error) {
	data := c
	if len(data) > 100 {
		data = c[:100]
	}

	idx := strings.Index(data, "\n")
	if idx != -1 {
		data = c[:idx]
	}

	matches := patterns.MessageLink.FindStringSubmatch(data)
	if len(matches) == 0 {
		return nil, c, nil
	}
	link := matches[0]
	guildID := matches[1]
	channelID := matches[2]
	messageID := matches[3]

	if guildID != vars.GuildID {
		return nil, "", errors.New("Replied message lives on different guild.")
	}
	if channelID != vars.ChannelID {
		return nil, "", errors.New("Replied message lives on different guild's channel.")
	}

	c = c[len(link):]
	if strings.TrimSpace(c) == "" {
		return nil, "", errEmptyMessage
	}

	ref := &discordgo.MessageReference{
		MessageID: messageID,
		GuildID:   vars.GuildID,
	}
	return ref, c, nil
}
