package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jozsefsallai/thimble-bot-telegram/commands"
	"github.com/jozsefsallai/thimble-bot-telegram/config"
	"github.com/jozsefsallai/thimble-bot-telegram/utils"

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

	bot.Handle("/start", func(m *tb.Message) {
		message := fmt.Sprintf("Hello there, %s! I am ready to serve you :)", m.Sender.Username)
		bot.Send(m.Chat, message)
	})

	bot.Handle("/8ball", commands.EightBallCommand(bot))
	utils.MultiCommand(bot, []string{"/cat", "/randomcat"}, commands.ShibeAPICommand(bot, "cat"))
	utils.MultiCommand(bot, []string{"/bird", "/birb", "/randombird"}, commands.ShibeAPICommand(bot, "bird"))
	bot.Handle("/h", commands.HCommand(bot))
	bot.Handle("/reverse", commands.ReverseCommand(bot))
	bot.Handle("/ship", commands.ShipCommand(bot))

	bot.Start()
}
