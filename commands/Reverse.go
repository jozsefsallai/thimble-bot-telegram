package commands

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func reverse(str string) string {
	n := len(str)
	runes := make([]rune, n)

	for _, character := range str {
		n--
		runes[n] = character
	}

	return string(runes[n:])
}

// ReverseCommand will return the reverse of a given string
func ReverseCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		bot.Send(m.Chat, reverse(m.Payload))
	}
}
