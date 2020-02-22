package commands

import (
	"github.com/jozsefsallai/go-derpi"

	tb "gopkg.in/tucnak/telebot.v2"
)

func PonyCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		if len(m.Payload) == 0 {
			bot.Send(m.Chat, "Please provide a search query.")
			return
		}

		client := derpi.Init()
		search := client.Search(derpi.SearchInitOpts{FilterID: 137238})

		image, err := search.RandomImage(m.Payload)
		if err != nil {
			bot.Send(m.Chat, "Failed to fetch the requested content.")
			return
		}

		if len(image.Format) == 0 {
			bot.Send(m.Chat, "No results found for your query.")
			return
		}

		bot.Send(m.Chat, image.Representations.Full)
	}
}
