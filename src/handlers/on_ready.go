package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(_ *discordgo.Session, r *discordgo.Ready) {
	log.Printf(
		"Logged in successfully: %s#%s (id=%s) (session=%s) (version=%d)\n",
		r.User.Username,
		r.User.Discriminator,
		r.User.ID,
		r.SessionID,
		r.Version,
	)
}
