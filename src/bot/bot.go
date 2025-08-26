package bot

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/4nonch/echochamber-dc/src/commands"
	"github.com/4nonch/echochamber-dc/src/handlers"
	"github.com/4nonch/echochamber-dc/src/vars"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

type Bot struct {
	session            *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
}

// Initializes bot instance and applies the required configurations
func NewBot() *Bot {
	session, err := discordgo.New("Bot " + vars.BotToken)
	if err != nil {
		log.Fatal("Failed to initialized bot:", err)
	}

	session.StateEnabled = true
	session.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentGuilds |
			discordgo.IntentDirectMessages |
			discordgo.IntentGuildMessages |
			discordgo.IntentGuildMembers,
	)

	handlers.Register(session)

	bot := &Bot{session: session}

	if vars.ProxyUrl != "" {
		bot.setupProxy(vars.ProxyUrl)
	}

	return bot
}

// Opens websocket connection to Discord
func (b *Bot) Start() {
	err := b.session.Open()
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}

	b.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands.All))
	for i, c := range commands.All {
		registered, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", c.Command)
		if err != nil {
			log.Panicf("Failed to register command \"%v\": %v", c.Command.Name, err)
		}
		b.registeredCommands[i] = registered
	}
}

// Closes websocket connection to Discord
func (b *Bot) Close() {
	log.Println("Shutting down...")
	b.session.Close()
	for _, c := range b.registeredCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", c.ID)
		if err != nil {
			log.Panicf("Failed to delete command \"%v\": %v", c.Name, err)
		}
	}
	b.registeredCommands = []*discordgo.ApplicationCommand{}
	log.Println("Shut down completed.")
}

// Instructs the current process to wait for future events from Discord.
// Handles the user's interrupt signal to stop processing when it occurs.
func (b *Bot) Await() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Received shutting signal.")
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
