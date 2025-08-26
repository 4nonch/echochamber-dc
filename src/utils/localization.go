package utils

import (
	"github.com/bwmarrin/discordgo"
)

type Localization struct {
	Loc discordgo.Locale
	Msg string
}

type LocaleMap map[discordgo.Locale]string

const DEFAULT_LOCALE discordgo.Locale = "_"

func GetLocalized(m LocaleMap, l discordgo.Locale) string {
	if value, ok := m[l]; ok {
		return value
	}
	return m[DEFAULT_LOCALE]
}

func MakeLocaleMap(original string, options ...*Localization) LocaleMap {
	localeMap := make(LocaleMap, len(options)+1)
	localeMap[DEFAULT_LOCALE] = original

	for _, loc := range options {
		localeMap[loc.Loc] = loc.Msg
	}

	return localeMap
}
