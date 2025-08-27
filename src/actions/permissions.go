package actions

import (
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/4nonch/echochamber-dc/src/utils"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

var (
	_tooBigContentMsg = utils.MakeLocaleMap(
		"You're message is too big (%d). Maximum allowed size: %d characters.",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Ваше сообщение слишком велико (%d). Максимально допустимый размер: %d символов.",
		},
	)
	_tooMuchAttachmentsMsg = utils.MakeLocaleMap(
		"Too much attachments (%d). Maximum allowed count: %d",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Слишком много прикреплённых файлов (%d). Максимально допустимое количество: %d файлов.",
		},
	)
	_tooBigAttachmentsSizeMsg = utils.MakeLocaleMap(
		"Attachments are too big (%.2f Mb). Maximum allowed size: %.2f Mb.",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Слишком большой вес файлов (%.2f Мб). Максимально допустимый размер: %.2f Мб.",
		},
	)
)

func CouldViewChannel(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	guild, err := s.State.Guild(vars.GuildID)
	if err == nil && guild.OwnerID == m.Author.ID {
		return true
	}

	perms, err := GetChannelPermissions(s, m.Author.ID)
	if err != nil {
		msg := fmt.Sprintf("Unable to get user's permissions for channel \"%s\": %v", vars.ChannelID, err)
		log.Printf(msg)
		SendMessage(msg, s, m)
		return false
	}

	return (perms&discordgo.PermissionViewChannel != 0) ||
		(perms&discordgo.PermissionAdministrator != 0)
}

func HasValidContent(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	count := utf8.RuneCountInString(m.Content)
	if count > vars.MaxMessageChars {
		tooBigContentMsg := utils.GetLocalized(_tooBigContentMsg, discordgo.Locale(m.Author.Locale))
		tooBigContentMsg = fmt.Sprintf(tooBigContentMsg, count, vars.MaxMessageChars)
		SendMessage(tooBigContentMsg, s, m)
		return false
	}
	return true
}

func HasValidMedia(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// Max attachments count check
	if len(m.Attachments) > vars.MaxAttachmentsCount {
		tooMuchAttachmentsMsg := utils.GetLocalized(_tooMuchAttachmentsMsg, discordgo.Locale(m.Author.Locale))
		tooMuchAttachmentsMsg = fmt.Sprintf(tooMuchAttachmentsMsg, len(m.Attachments), vars.MaxAttachmentsCount)
		SendMessage(tooMuchAttachmentsMsg, s, m)
		return false
	}

	// Max attachments size check
	var total int
	for _, a := range m.Attachments {
		total += a.Size
	}
	if total > vars.MaxAttachmentsBytes {
		tooBigAttachmentSizeMsg := utils.GetLocalized(
			_tooBigAttachmentsSizeMsg,
			discordgo.Locale(m.Author.Locale),
		)
		tooBigAttachmentSizeMsg = fmt.Sprintf(
			tooBigAttachmentSizeMsg,
			float32(total)/1024/1024,
			float32(vars.MaxAttachmentsBytes)/1024/1024,
		)
		SendMessage(tooBigAttachmentSizeMsg, s, m)
		return false
	}
	return true
}
