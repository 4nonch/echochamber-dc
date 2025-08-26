package vars

import (
	"log"
	"os"
)

var (
	BotToken  = getEnv("BOT_TOKEN")
	GuildID   = getEnv("GUILD_ID")
	ChannelID = getEnv("CHANNEL_ID")
	ProxyUrl  = getEnv("PROXY_URL", "")
)

func getEnv(lookup string, defaultValue ...string) string {
	env, exists := os.LookupEnv(lookup)
	if exists {
		return env
	}
	if len(defaultValue) != 0 {
		return defaultValue[0]
	}
	log.Fatalf("Unable to lookup ENV \"%s\", execution impossible", lookup)
	return ""
}
