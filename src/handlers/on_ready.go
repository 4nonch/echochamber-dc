package handlers

import (
	"log"

	"github.com/4nonch/echochamber-dc/src/cache"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.add(OnReady)
}

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf(
		"Logged in successfully: %s#%s (id=%s) (session=%s) (version=%d)\n",
		r.User.Username,
		r.User.Discriminator,
		r.User.ID,
		r.SessionID,
		r.Version,
	)

	emojis, err := s.GuildEmojis(vars.GuildID)
	if err != nil {
		log.Fatalf(
			"Failed to initialize: unable to retrieve emoji data from target's guild %s. No access?",
			vars.GuildID,
		)
	}
	cache.Emojis.Set(emojis)
}
