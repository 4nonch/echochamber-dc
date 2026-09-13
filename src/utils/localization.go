// Package utils provides various helpers
package utils

import (
	"github.com/bwmarrin/discordgo"
)

type Localization struct {
	Loc discordgo.Locale
	Msg string
}

type LocaleMap map[discordgo.Locale]string

const DefaultLocale discordgo.Locale = "_"

func GetLocalized(m LocaleMap, l discordgo.Locale) string {
	if value, ok := m[l]; ok {
		return value
	}
	return m[DefaultLocale]
}

func MakeLocaleMap(original string, options ...*Localization) LocaleMap {
	localeMap := make(LocaleMap, len(options)+1)
	localeMap[DefaultLocale] = original

	for _, loc := range options {
		localeMap[loc.Loc] = loc.Msg
	}

	return localeMap
}
