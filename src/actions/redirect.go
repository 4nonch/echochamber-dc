package actions

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

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

	if len(m.Attachments) == 0 {
		_, err := s.ChannelMessageSend(vars.ChannelID, m.Content)
		if err != nil {
			log.Printf("Failed to redirect message \"%s\": %v", m.Content, err)
		} else {
			sendSuccess(s, m)
		}
		return
	}

	streamed, errs := concurrentDownload(m.Attachments)
	if len(errs) != 0 {
		for err := range errs {
			log.Println("Failed to get attachment: ", err)
		}
		SendMessage("An error occurred while trying to download attachments.", s, m)
		return
	}

	msg := &discordgo.MessageSend{
		Content: m.Content,
		Files:   streamed.files,
	}
	_, err := s.ChannelMessageSendComplex(vars.ChannelID, msg)
	streamed.close()
	if err != nil {
		errMsg := fmt.Sprintf(
			"Failed to redirect message \"%s\" (attachments: %d): %v",
			m.Content,
			len(m.Attachments),
			err,
		)
		log.Println(errMsg)
		SendMessage(errMsg, s, m)
	}
	sendSuccess(s, m)
}

func sendSuccess(s *discordgo.Session, m *discordgo.MessageCreate) {
	SendMessage(utils.GetLocalized(_successMsg, discordgo.Locale(m.Author.Locale)), s, m)
}

func concurrentDownload(attachments []*discordgo.MessageAttachment) (streamedFiles, chan error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(attachments))
	streamed := streamedFiles{
		files: make([]*discordgo.File, len(attachments)),
		resps: make([]io.ReadCloser, len(attachments)),
	}

	for i, a := range attachments {
		wg.Add(1)
		go func(i int, a *discordgo.MessageAttachment) {
			defer wg.Done()
			res, err := vars.Client.Get(a.URL)

			if err != nil {
				errChan <- err
				return
			}
			if res.StatusCode != http.StatusOK {
				res.Body.Close()
				errChan <- errors.New(fmt.Sprintf("Failed to fetch %s, status code: %d", a.URL, res.StatusCode))
				return
			}

			streamed.files[i] = &discordgo.File{
				Name:        a.Filename,
				ContentType: a.ContentType,
				Reader:      res.Body,
			}
			streamed.resps[i] = res.Body
		}(i, a)
	}

	wg.Wait()
	return streamed, errChan
}

type streamedFiles struct {
	files []*discordgo.File
	resps []io.ReadCloser
}

func (f *streamedFiles) close() {
	for _, file := range f.resps {
		file.Close()
	}
}
