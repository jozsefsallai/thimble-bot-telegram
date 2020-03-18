package utils

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// PhotoHandler is a function that returns a closure containing a Telegram message
type PhotoHandler func(*tb.Bot) func(*tb.Message)

func getCommandFromCaption(caption string) string {
	return strings.Fields(caption)[0][1:]
}

// HandlePhotos is used for photo-related commands (where the commands are) in
// the caption of the photo
func HandlePhotos(bot *tb.Bot, handlers map[string]PhotoHandler) interface{} {
	return func(m *tb.Message) {
		if len(m.Caption) == 0 {
			return
		}

		command := getCommandFromCaption(m.Caption)
		handler := handlers[command]

		if handler == nil {
			return
		}

		handler(bot)(m)
	}
}
