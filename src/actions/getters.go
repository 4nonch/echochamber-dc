package actions

import (
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
