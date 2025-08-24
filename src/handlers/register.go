package handlers

import "github.com/bwmarrin/discordgo"

func Register(s *discordgo.Session) {
	for _, h := range []any{
		OnReady,
		OnMessageCreate,
		OnInteractionCreate,
	} {
		s.AddHandler(h)
	}
}
