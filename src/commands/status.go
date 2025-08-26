package commands

import (
	"fmt"
	"log"

	"github.com/4nonch/echochamber-dc/src/actions"
	"github.com/4nonch/echochamber-dc/src/utils"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

var (
	StatusCommand = &Command{
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
	_botInfoMsg = utils.MakeLocaleMap(
		"Target guild: %s\n"+
			"Target channel: %s",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Целевой сервер: %s\n" +
				"Целевой канал: %s",
		},
	)
	_userAllowedMsg = utils.MakeLocaleMap(
		"You're allowed to write messages in target's guild channel.",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Вам разрешено отправлять сообщения в целевой канал сервера.",
		},
	)
	_userNotAllowedMsg = utils.MakeLocaleMap(
		"You're not a member of target's guild channel. Can't forward messages.",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Вы не являетесь участником целевого канала сервера. Отправка сообщений запрещена.",
		},
	)
)

func onStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := actions.GetGuild(s)
	if err != nil {
		message := fmt.Sprintf("Unable to query guild \"%s\": %v", vars.GuildID, err)
		log.Printf(message)
		actions.SendInteractionMessage(message, s, i)
		return
	}

	channel, err := actions.GetChannel(s)
	if err != nil {
		message := fmt.Sprintf("Unable to query channel \"%s\": %v", vars.ChannelID, err)
		log.Printf(message)
		actions.SendInteractionMessage(message, s, i)
		return
	}

	member, err := actions.GetGuildMember(s, i.User.ID)
	if err != nil {
		message := fmt.Sprintf("Unable to get user's status, guild \"%s\": %v", vars.ChannelID, err)
		log.Printf(message)
		actions.SendInteractionMessage(message, s, i)
		return
	}

	msg := _makeStatusResponse(guild, channel, member, i)
	actions.SendInteractionMessage(msg, s, i)
}

func _makeStatusResponse(
	g *discordgo.Guild,
	c *discordgo.Channel,
	m *discordgo.Member,
	i *discordgo.InteractionCreate,
) string {
	botInfoMsg := fmt.Sprintf(
		utils.GetLocalized(_botInfoMsg, i.Locale),
		g.Name, c.Name,
	)

	var permissionMsg string
	if m != nil {
		permissionMsg = utils.GetLocalized(_userAllowedMsg, i.Locale)
	} else {
		permissionMsg = utils.GetLocalized(_userNotAllowedMsg, i.Locale)
	}

	return fmt.Sprintf("%s\n\n%s", botInfoMsg, permissionMsg)
}
