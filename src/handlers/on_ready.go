package handlers

import (
	"log"
	"strings"

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

	guilds := make([]string, len(r.Guilds))
	for i, g := range r.Guilds {
		guilds[i] = g.Name
	}
	if guilds[0] != "" {
		log.Print("Available guilds: ", strings.Join(guilds, ","))
	}
}
