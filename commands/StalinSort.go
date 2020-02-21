package commands

import (
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func isNumeric(n string) (int, bool) {
	number, err := strconv.ParseInt(n, 0, 64)
	return int(number), err == nil
}

func filterInts(arr []string) []int {
	var result []int

	for _, entry := range arr {
		number, ok := isNumeric(entry)
		if ok {
			result = append(result, number)
		}
	}

	return result
}

func stalinSort(arr []int) []int {
	MinInt := -(int(^uint(0) >> 1)) - 1
	var result []int

	for _, current := range arr {
		if current >= MinInt {
			MinInt = current
			result = append(result, current)
		}
	}

	return result
}

func joinResults(arr []int) string {
	result := ""

	for idx, current := range arr {
		if idx == 0 {
			result += strconv.Itoa(current)
		} else {
			result += fmt.Sprintf(" %s", strconv.Itoa(current))
		}
	}

	return result
}

// StalinSortCommand takes the numbers from the user input and sorts them
// using the O(n) Stalin Sort algorithm
func StalinSortCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		numbers := filterInts(strings.Split(m.Payload, " "))
		if len(numbers) == 0 {
			bot.Send(m.Chat, "Please provide numbers to sort.")
			return
		}

		sorted := stalinSort(numbers)

		response := fmt.Sprintf("Here it is sorted: `%s`", joinResults(sorted))
		bot.Send(m.Chat, response, tb.ModeMarkdown)
	}
}
