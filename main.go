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
	bot.Handle("/addnana", commands.AddNanaCommand(bot))
	bot.Handle("/choice", commands.ChoiceCommand(bot))
	bot.Handle("/delnana", commands.DeleteNanaCommand(bot))
	bot.Handle("/delfapu", commands.DeleteFapuCommand(bot))
	utils.MultiCommand(bot, aliases.For["Faputa"], commands.FaputaCommand(bot))
	bot.Handle("/flip", commands.FlipCommand(bot))
	utils.MultiCommand(bot, aliases.For["Nanachi"], commands.NanachiCommand(bot))
	utils.MultiCommand(bot, aliases.For["RandomCat"], commands.ShibeAPICommand(bot, "cat"))
	utils.MultiCommand(bot, aliases.For["RandomBird"], commands.ShibeAPICommand(bot, "bird"))
	utils.MultiCommand(bot, aliases.For["RandomBunny"], commands.BunnyCommand(bot))
	utils.MultiCommand(bot, aliases.For["RandomDog"], commands.DogCommand(bot))
	bot.Handle("/h", commands.HCommand(bot))
	bot.Handle("/reverse", commands.ReverseCommand(bot))
	utils.MultiCommand(bot, aliases.For["RandomPony"], commands.PonyCommand(bot))
	bot.Handle("/setsource", commands.SetSourceCommand(bot))
	bot.Handle("/ship", commands.ShipCommand(bot))
	bot.Handle("/slap", commands.SlapCommand(bot))
	bot.Handle("/stalinsort", commands.StalinSortCommand(bot))
	utils.MultiCommand(bot, aliases.For["Strawpoll"], commands.StrawpollCommand(bot, false))
	utils.MultiCommand(bot, aliases.For["StrawpollMulti"], commands.StrawpollCommand(bot, true))

	photoHandlers := make(map[string]utils.PhotoHandler)
	photoHandlers["addnana"] = commands.AddNanaCommand
	photoHandlers["addfapu"] = commands.AddFapuCommand

	bot.Handle(tb.OnPhoto, utils.HandlePhotos(bot, photoHandlers))

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
