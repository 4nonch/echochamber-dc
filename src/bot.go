package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/4nonch/echochamber-dc/src/commands"
	"github.com/4nonch/echochamber-dc/src/handlers"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

type Bot struct {
	session *discordgo.Session
}

// Initializes bot instance and applies the required configurations
func NewBot(config *Config) *Bot {
	session, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal("Failed to initialized bot:", err)
	}

	session.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentDirectMessages |
			discordgo.IntentGuildMessages |
			discordgo.IntentGuildMembers,
	)

	handlers.Register(session)

	bot := &Bot{session: session}

	if config.ProxyUrl != "" {
		bot.setupProxy(config.ProxyUrl)
	}

	return bot
}

// Opens websocket connection to Discord
func (b *Bot) Start() {
	err := b.session.Open()
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}

	for _, c := range commands.All {
		_, err = b.session.ApplicationCommandCreate(b.session.State.User.ID, "", c.Command)
		if err != nil {
			log.Panicf("Failed to register command \"%v\": %v", c.Command.Name, err)
		}
	}
}

// Closes websocket connection to Discord
func (b *Bot) Close() {
	b.session.Close()
	for _, c := range commands.All {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", c.Command.ID)
		if err != nil {
			log.Panicf("Failed to delete command \"%v\": %v", c.Command.Name, err)
		}
	}
}

// Instructs the current process to wait for future events from Discord.
// Handles the user's interrupt signal to stop processing when it occurs.
func (b *Bot) Await() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

// Configures discordgo to use user-specified proxy URL for it's connections
// (useful for circumventing Discord blocks in third-world countries)
func (b *Bot) setupProxy(proxyUrl string) {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		log.Fatalf(
			"Failed to parse PROXY_URL \"%s\": %s",
			proxyUrl,
			err,
		)
	}
	proxyFactory := http.ProxyURL(proxy)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxyFactory,
		},
	}
	dialer := &websocket.Dialer{
		Proxy: proxyFactory,
	}

	b.session.Client = client
	b.session.Dialer = dialer
}
