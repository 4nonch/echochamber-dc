package services

import (
	"errors"
	"strings"

	"github.com/4nonch/echochamber-dc/src/patterns"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

// Return MessageReference if original content contained a reply.
// Link to the replied message will be removed from original content
func extractReference(m *discordgo.MessageCreate) (*discordgo.MessageReference, string, error) {
	c := m.Content
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

	ref := &discordgo.MessageReference{
		MessageID: messageID,
		GuildID:   vars.GuildID,
	}
	return ref, c, nil
}
