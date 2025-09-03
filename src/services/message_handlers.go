package services

import (
	"errors"
	"strings"

	"github.com/4nonch/echochamber-dc/src/cache"
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

type _emojiRep struct {
	left, right int
	replace     string
}

// Parses string in search of guild's emojis and formats the content with proper emoji codes,
// so the emojis could be properly displayed in target's channel.
// Will also replace :emoji: wrapped in “ formatting symbols, so it could be like `<:emoji:123>` in the end
func formatGuildEmojis(content string) string {
	if len(content) < 3 {
		return content
	}

	first := strings.Index(content, ":")
	if first == -1 || first >= len(content)-2 {
		return content
	}

	var code string
	var b strings.Builder
	replacements := make([]_emojiRep, 0, vars.MaxMessageChars/3)
	b.Grow(vars.MaxMessageChars)

	left := first
	for right := first + 1; right < len(content); right++ {
		if content[right] != ':' {
			continue
		}
		code = cache.Emojis.GetCode(content[left+1 : right])
		if code == "" {
			left = right
			continue
		}
		replacements = append(replacements, _emojiRep{left, right, code})
	}

	lastIdx := 0
	for _, node := range replacements {
		b.WriteString(content[lastIdx:node.left])
		b.WriteString(node.replace)
		lastIdx = node.right + 1
	}
	b.WriteString(content[lastIdx:])

	return b.String()
}
