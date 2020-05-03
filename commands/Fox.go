package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tb "gopkg.in/tucnak/telebot.v2"
)

// RandomFoxAPIResponse contains the JSON response returned by randomfox.ca
type RandomFoxAPIResponse struct {
	Image string
	Link  string
}

func fail(bot *tb.Bot, m *tb.Message) {
	bot.Send(m.Chat, "Failed to fetch the picture.")
}

func FoxCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		apiURL := "https://randomfox.ca/floof"
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fail(bot, m)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fail(bot, m)
			return
		}

		defer res.Body.Close()

		var result RandomFoxAPIResponse

		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &result)

		image := &tb.Photo{File: tb.FromURL(result.Image)}
		bot.Send(m.Chat, image)
	}
}
