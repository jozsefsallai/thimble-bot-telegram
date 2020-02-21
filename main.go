package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jozsefsallai/thimble-bot-telegram/aliases"
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

	bot.Handle("/help", commands.HelpCommand(bot))

	utils.MultiCommand(bot, aliases.For["8ball"], commands.EightBallCommand(bot))
	utils.MultiCommand(bot, aliases.For["RandomCat"], commands.ShibeAPICommand(bot, "cat"))
	utils.MultiCommand(bot, aliases.For["RandomBird"], commands.ShibeAPICommand(bot, "bird"))
	utils.MultiCommand(bot, aliases.For["RandomBunny"], commands.BunnyCommand(bot))
	bot.Handle("/h", commands.HCommand(bot))
	bot.Handle("/reverse", commands.ReverseCommand(bot))
	bot.Handle("/ship", commands.ShipCommand(bot))
	bot.Handle("/stalinsort", commands.StalinSortCommand(bot))

	fmt.Println(strings.Join([]string{
		"  _______ _     _           _     _        ____        _   ",
		" |__   __| |   (_)         | |   | |      |  _ \\      | |  ",
		"    | |  | |__  _ _ __ ___ | |__ | | ___  | |_) | ___ | |_ ",
		"    | |  | '_ \\| | '_ ` _ \\| '_ \\| |/ _ \\ |  _ < / _ \\| __|",
		"    | |  | | | | | | | | | | |_) | |  __/ | |_) | (_) | |_ ",
		"    |_|  |_| |_|_|_| |_| |_|_.__/|_|\\___| |____/ \\___/ \\__|",
	}, "\n"))
	fmt.Println("\nApp started successfully...")

	bot.Start()
}
