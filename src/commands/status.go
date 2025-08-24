package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var StatusCommand = &Command{
	Command: &discordgo.ApplicationCommand{
		Name:        "status",
		Description: "Print out bot configurations and user's status",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.Russian: "статус",
		},
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.Russian: "Выводит настройки бота и статус пользователя",
		},
		Contexts: &[]discordgo.InteractionContextType{
			discordgo.InteractionContextBotDM,
			discordgo.InteractionContextPrivateChannel,
		},
	},
	Handler: onStatusCommand,
}

func onStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("\"Status\" Command Triggered")
}
