package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tb "gopkg.in/tucnak/telebot.v2"
)

// ShibeAPICommand will fetch a random picture from shibe.online
func ShibeAPICommand(bot *tb.Bot, mode string) interface{} {
	return func(m *tb.Message) {
		var endpoint string

		switch mode {
		case "cat":
			endpoint = "cats"
		case "bird":
			endpoint = "birds"
		default:
			endpoint = "cats"
		}

		apiURL := "https://shibe.online/api/" + endpoint
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			bot.Send(m.Chat, "Failed to fetch the picture.")
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			bot.Send(m.Chat, "Failed to fetch the picture.")
		}

		defer res.Body.Close()

		var results []string

		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &results)

		image := &tb.Photo{File: tb.FromURL(results[0])}
		bot.Send(m.Chat, image)
	}
}
