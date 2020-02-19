package commands

import (
	"fmt"
	"strings"

	"github.com/jozsefsallai/thimble-bot-telegram/aliases"

	tb "gopkg.in/tucnak/telebot.v2"
)

type CommandDetails struct {
	NameWithParams string
	Aliases        []string
	Description    string
}

func getCommandRow(command CommandDetails) string {
	var name string
	var description string

	if len(command.Aliases) > 0 {
		name = strings.Join(command.Aliases, ", ")
	}

	if command.NameWithParams != "" {
		name = fmt.Sprintf("/%s", command.NameWithParams)
	}

	description = command.Description

	return fmt.Sprintf("*%s*\n%s\n", name, description)
}

// HelpCommand will return the available commands with a description
func HelpCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		commands := []CommandDetails{
			{
				NameWithParams: "8ball [question]",
				Description:    "Ask the 8-Ball any question.",
			},
			{
				Aliases:     aliases.For["RandomCat"],
				Description: "Get a random picture, GIF, or video of a cat.",
			},
			{
				Aliases:     aliases.For["RandomBird"],
				Description: "Get a random picture, GIF, or video of a bird.",
			},
			{
				Aliases:     aliases.For["RandomBunny"],
				Description: "Get a random GIF of a bunny",
			},
			{
				NameWithParams: "h",
				Description:    "h",
			},
			{
				NameWithParams: "reverse [input]",
				Description:    "Get the reverse of a given string.",
			},
			{
				NameWithParams: "ship [person1] x|and [person2]",
				Description:    "Get the love compatibility of two people. The values must be separated by \" x \" or \" and \".",
			},
		}

		var rows []string

		for _, command := range commands {
			rows = append(rows, getCommandRow(command))
		}

		bot.Send(m.Chat, strings.Join(rows, "\n"), tb.ModeMarkdown)
	}
}
