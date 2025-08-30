package actions

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
)

func GetGuild(s *discordgo.Session) (*discordgo.Guild, error) {
	guild, err := s.State.Guild(vars.GuildID)
	if err == nil {
		return guild, nil
	}
	guild, err = s.Guild(vars.GuildID)
	return guild, err
}

func GetChannel(s *discordgo.Session) (*discordgo.Channel, error) {
	channel, err := s.State.Channel(vars.ChannelID)
	if err == nil {
		return channel, nil
	}
	channel, err = s.Channel(vars.ChannelID)
	return channel, err
}

func GetGuildMember(s *discordgo.Session, userID string) (*discordgo.Member, error) {
	member, err := s.State.Member(vars.GuildID, userID)
	if err == nil {
		return member, nil
	}
	member, err = s.GuildMember(vars.GuildID, userID)
	return member, err
}

func GetChannelPermissions(s *discordgo.Session, userID string) (int64, error) {
	perms, err := s.State.UserChannelPermissions(userID, vars.ChannelID)
	if err == nil {
		return perms, nil
	}
	perms, err = s.UserChannelPermissions(userID, vars.ChannelID)
	return perms, err
}

// Used to fetch attachment bodies.
// If success, don't forget to close the body after use (StreamFiles.Close)
func GetAttachments(attachments []*discordgo.MessageAttachment) (StreamFiles, chan error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(attachments))
	stream := StreamFiles{
		Files: make([]*discordgo.File, len(attachments)),
		Resps: make([]io.ReadCloser, len(attachments)),
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

			stream.Files[i] = &discordgo.File{
				Name:        a.Filename,
				ContentType: a.ContentType,
				Reader:      res.Body,
			}
			stream.Resps[i] = res.Body
		}(i, a)
	}

	wg.Wait()
	return stream, errChan
}
