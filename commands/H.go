package commands

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

// HCommand will just say "h"
func HCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		bot.Send(m.Sender, "h")
	}
}
