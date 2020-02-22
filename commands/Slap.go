package commands

import (
	"fmt"
	"math/rand"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SlapCommand for Romana
func SlapCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		if m.Sender.ID != 805550175 {
			bot.Send(m.Chat, "You are not Romana")
			return
		}

		if len(m.Payload) == 0 {
			bot.Send(m.Chat, "You didn't tell me who to slap")
			return
		}

		endings := []string{
			"That is not very nice...",
			"I mean they kinda deserved it.",
			"Now they are hurt.",
			"Press F to pay respects.",
			"I'm gonna sue.",
			"They promised revenge.",
		}

		ending := endings[rand.Intn(len(endings))]

		bot.Send(m.Chat, fmt.Sprintf("👏 @%s slapped %s. %s", m.Sender.Username, m.Payload, ending))
	}
}
