package commands

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"

	tb "gopkg.in/tucnak/telebot.v2"
)

// EightBallCommand will return a seeded answer based on the user's input
func EightBallCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		if len(m.Payload) == 0 {
			bot.Send(m.Chat, "Please ask a valid question.")
			return
		}

		answers := [...]string{
			"Certainly so!",
			"Undoubtedly yes!",
			"Definitely.",
			"Very much yes!",
			"Yep!",
			"Most likely, yeah!",
			"Maybe-maybe...",
			"Probably.",
			"I cannot tell you that right now.",
			"Negative.",
			"I think the answer is no.",
			"Not really.",
			"Nope.",
			"No.",
			"Very doubtful.",
		}

		hasher := sha256.New()
		io.WriteString(hasher, m.Payload)

		var seed uint64 = binary.BigEndian.Uint64(hasher.Sum(nil))
		rand.Seed(int64(seed))

		answer := answers[rand.Intn(len(answers))]
		output := fmt.Sprintf("*🎱 The Great 8-Ball Says...*\n%s", answer)

		bot.Send(m.Chat, output, tb.ModeMarkdown)
		return
	}
}
