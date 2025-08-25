package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Command *discordgo.ApplicationCommand
	Handler func(*discordgo.Session, *discordgo.InteractionCreate)
}

var (
	All = []*Command{
		StatusCommand,
	}
	Handlers = allCommandHandlers()
)

func allCommandHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	commands := make(
		map[string]func(*discordgo.Session, *discordgo.InteractionCreate),
		len(All),
	)
	for _, c := range All {
		commands[c.Command.Name] = c.Handler
	}
	return commands
}
