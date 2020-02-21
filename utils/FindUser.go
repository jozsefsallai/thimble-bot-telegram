package utils

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

// FindUser will try to find a tb.MessageEntity containing a mention
func FindUser(entities []tb.MessageEntity) tb.MessageEntity {
	for _, entity := range entities {
		if entity.Type == tb.EntityTMention || entity.Type == tb.EntityMention {
			return entity
		}
	}

	return tb.MessageEntity{Length: 0}
}
