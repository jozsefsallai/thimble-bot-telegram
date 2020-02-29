package commands

import (
	"strings"

	"github.com/jozsefsallai/thimble-bot-telegram/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// ChoiceCommand will pick a random value from an array
func ChoiceCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		if len(m.Payload) == 0 {
			bot.Send(m.Chat, "The command expects parameters separated by pipe characters (|).")
			return
		}

		choices := strings.Split(m.Payload, "|")
		filteredChoices := utils.RemoveEmptyStrings(choices)

		if len(filteredChoices) == 0 {
			bot.Send(m.Chat, "Please provide things to pick from, delimited by | characters.")
			return
		}

		result := strings.TrimSpace(utils.ChoiceString(filteredChoices))
		bot.Send(m.Chat, result)
	}
}
