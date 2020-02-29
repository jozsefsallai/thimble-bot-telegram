package commands

import (
	"fmt"
	"strings"

	"github.com/jozsefsallai/go-strawpoll"
	"github.com/jozsefsallai/thimble-bot-telegram/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

// StrawpollCommand will create a poll on strawpoll.me
func StrawpollCommand(bot *tb.Bot, multi bool) interface{} {
	return func(m *tb.Message) {
		if len(m.Payload) == 0 {
			bot.Send(
				m.Chat,
				"The command expects the following input: `Poll Title: choice 1 | choice 2.`",
				tb.ModeMarkdown,
			)
			return
		}

		input := strings.SplitN(m.Payload, ":", 2)
		filteredInput := utils.RemoveEmptyStrings(input)

		if len(filteredInput) != 2 {
			bot.Send(
				m.Chat,
				"The correct syntax is: `/strawpoll Title: option 1 | option 2 | option 3`",
				tb.ModeMarkdown,
			)
			return
		}

		title := filteredInput[0]
		choices := strings.Split(filteredInput[1], "|")
		filteredChoices := utils.RemoveEmptyStrings(choices)

		if len(filteredChoices) < 2 {
			bot.Send(m.Chat, "You need to provide at least 2 choices.")
			return
		}

		poll, err := strawpoll.Create(
			title,
			filteredChoices,
			multi,
			strawpoll.DupcheckNormal,
			false,
		)

		if err != nil {
			bot.Send(m.Chat, "Failed to create the strawpoll")
			return
		}

		bot.Send(
			m.Chat,
			fmt.Sprintf("*%s*\nhttps://www.strawpoll.me/%d", poll.Title, poll.ID),
			tb.ModeMarkdown,
			tb.NoPreview,
		)
	}
}
