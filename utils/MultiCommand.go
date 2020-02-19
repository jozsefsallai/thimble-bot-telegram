package utils

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

// MultiCommand is a slim wrapper that allows you to specify the same callback
// for multiple commands/handlers
func MultiCommand(bot *tb.Bot, commands []string, callback interface{}) {
	for _, command := range commands {
		bot.Handle(command, callback)
	}
}
