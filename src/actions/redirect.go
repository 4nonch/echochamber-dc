package actions

import (
	"fmt"
	"log"

	"github.com/4nonch/echochamber-dc/src/utils"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

var (
	_successMsg = utils.MakeLocaleMap(
		"Message successfully delivered.",
		&utils.Localization{
			Loc: discordgo.Russian,
			Msg: "Сообщение успешно доставлено.",
		},
	)
)

func RedirectMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !HasValidContent(s, m) || !HasValidMedia(s, m) {
		return
	}

	stream, failed := extractAttachments(s, m)
	if failed {
		return
	}

	dm := &discordgo.MessageSend{}
	dm.Content = m.Content
	dm.Files = stream.Files
	_, err := s.ChannelMessageSendComplex(vars.ChannelID, dm)
	stream.Close()

	if err != nil {
		msg := fmt.Sprintf(
			"Failed to redirect message \"%s\" (attachments: %d): %v",
			m.Content,
			len(m.Attachments),
			err,
		)
		log.Println(msg)
		SendMessage(msg, s, m)
		return
	}

	SendMessage(utils.GetLocalized(_successMsg, discordgo.Locale(m.Author.Locale)), s, m)
}

// Returns attachment files (if there is no one, Files attribute will be simply an empty slice).
// If second error is true - then error occurred and notified
func extractAttachments(s *discordgo.Session, m *discordgo.MessageCreate) (StreamFiles, bool) {
	var stream StreamFiles

	if len(m.Attachments) == 0 {
		return stream, false
	}

	var errs chan error
	stream, errs = GetAttachments(m.Attachments)
	if len(errs) != 0 {
		for err := range errs {
			log.Println("Failed to get attachment: ", err)
		}
		SendMessage("An error occurred while trying to download attachments.", s, m)
		return stream, true
	}
	return stream, false
}
