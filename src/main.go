package main

func main() {
	config := LoadConfig()
	bot := NewBot(config)
	bot.Start()
	bot.Await()
	bot.Close()
}
