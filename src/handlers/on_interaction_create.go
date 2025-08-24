package handlers

import (
	"log"

	"github.com/4nonch/echochamber-dc/src/commands"
	"github.com/bwmarrin/discordgo"
)

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand ||
		i.GuildID != "" {
		return
	}

	name := i.ApplicationCommandData().Name
	log.Println("Command triggered: ", name)

	if h, ok := commands.Handlers[name]; ok {
		h(s, i)
	}
}
