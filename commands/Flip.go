package commands

import (
	"github.com/jozsefsallai/thimble-bot-telegram/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// FlipCommand will flip a coin
func FlipCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		bot.Send(m.Chat, utils.ChoiceString([]string{"head", "tails"}))
	}
}
