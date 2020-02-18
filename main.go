package main

import (
	"log"
	"time"

	"github.com/jozsefsallai/thimble-bot-telegram/commands"
	"github.com/jozsefsallai/thimble-bot-telegram/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	conf := config.GetConfig()

	bot, err := tb.NewBot(tb.Settings{
		Token:  conf.Bot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/h", commands.HCommand(bot))
	bot.Handle("/reverse", commands.ReverseCommand(bot))
	bot.Handle("/ship", commands.ShipCommand(bot))

	bot.Start()
}
