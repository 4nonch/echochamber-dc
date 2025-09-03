package handlers

import (
	"github.com/4nonch/echochamber-dc/src/cache"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.add(OnGuildEmojisUpdate)
}

func OnGuildEmojisUpdate(s *discordgo.Session, ge *discordgo.GuildEmojisUpdate) {
	if ge.GuildID != vars.GuildID {
		return
	}
	cache.Emojis.Set(ge.Emojis)
}
