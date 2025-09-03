package cache

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// We'll keep cache of guild's emojis and will update them at GuildEmojisUpdate/Ready events
var Emojis = guildEmojis{}

type guildEmojis struct {
	mu   sync.RWMutex
	data map[string]string
}

func (e *guildEmojis) Set(emojis []*discordgo.Emoji) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data = make(map[string]string, len(emojis))
	for _, emoji := range emojis {
		e.data[emoji.Name] = emoji.ID
	}
}

func (e *guildEmojis) GetCode(key string) string {
	id, ok := e.data[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("<:%s:%s>", key, id)
}
