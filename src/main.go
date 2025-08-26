package main

import (
	"github.com/4nonch/echochamber-dc/src/bot"
)

func main() {
	b := bot.NewBot()
	b.Start()
	b.Await()
	b.Close()
}
