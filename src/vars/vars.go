package vars

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

var (
	// Main variables
	BotToken  = getEnv("BOT_TOKEN")
	GuildID   = getEnv("GUILD_ID")
	ChannelID = getEnv("CHANNEL_ID")

	// Discord's specific settings
	MaxMessageChars     = getEnvInt("MAX_MESSAGE_CHARS", 1800)
	MaxAttachmentsBytes = getEnvInt("MAX_ATTACHMENTS_BYTES", 10485760)
	MaxAttachmentsCount = getEnvInt("MAX_ATTACHMENTS_COUNT", 10)
)

var (
	// Proxy, if you'd need it
	proxyUrl = parseUrl(getEnv("PROXY_URL", ""))
	// Default HTTP client
	Client = getHttpClient()
	// Default websocket dialer
	Dialer = getWebsocketDialer()
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

func getEnvInt(lookup string, defaultValue ...int) int {
	var env string
	if len(defaultValue) != 0 {
		env = getEnv(lookup, strconv.Itoa(defaultValue[0]))
	} else {
		env = getEnv(lookup)
	}
	value, err := strconv.Atoi(env)
	if err != nil {
		log.Fatalf(
			"Environment variable \"%s\" (value: %v) can't be parsed as integer: %v",
			lookup, value, err,
		)
	}
	return value
}

func parseUrl(u string) *url.URL {
	if u == "" {
		return nil
	}
	r, err := url.Parse(u)
	if err != nil {
		log.Fatalf(
			"Failed to parse proxy url \"%s\": %s",
			u,
			err,
		)
	}
	return r
}

func getHttpClient() *http.Client {
	if proxyUrl == nil {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
}

func getWebsocketDialer() *websocket.Dialer {
	if proxyUrl == nil {
		return websocket.DefaultDialer
	}
	return &websocket.Dialer{
		Proxy: http.ProxyURL(proxyUrl),
	}
}
