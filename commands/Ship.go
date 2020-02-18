package commands

import (
	"fmt"
	"regexp"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func getDelimiter(input string) string {
	r, _ := regexp.Compile(" x | and ")
	return r.FindString(input)
}

func getASCIISum(subject string) int {
	sum := 0

	for _, character := range subject {
		sum += int(character)
	}

	return sum
}

func calculate(subjects []string) int {
	shipCount := 0

	for _, subject := range subjects {
		subjectLowercase := strings.ToLower(strings.TrimSpace(subject))
		shipCount += getASCIISum(subjectLowercase)
	}

	return shipCount % 100
}

// ShipCommand will calculate the love compatibility of two people
func ShipCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		delimiter := getDelimiter(m.Payload)
		if len(delimiter) == 0 {
			bot.Send(m.Sender, "The arguments you have provided are invalid. The correct form is `person1 x person2` or `person1 and person2`.")
			return
		}

		subjects := strings.Split(m.Payload, delimiter)
		if len(subjects) != 2 {
			bot.Send(m.Sender, "Too many arguments provided. Please enter only two.")
			return
		}

		compatibility := calculate(subjects)

		output := fmt.Sprintf("The ship compatibility of %s %s %s is %d%%.", subjects[0], strings.TrimSpace(delimiter), subjects[1], compatibility)
		bot.Send(m.Sender, output)
	}
}
