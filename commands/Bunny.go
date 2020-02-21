package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tb "gopkg.in/tucnak/telebot.v2"
)

// BunniesIOResponseGIFMedia is a structure that contains the GIF
// returned by api.bunnies.io
type BunniesIOResponseGIFMedia struct {
	Gif    string `json:"gif"`
	Poster string `json:"poster"`
}

// BunniesIOResponse is a structure that contains the full response
// returned by api.bunnies.io
type BunniesIOResponse struct {
	ID          int                       `json:"id"`
	Media       BunniesIOResponseGIFMedia `json:"media"`
	Source      string                    `json:"source"`
	ThisServed  int                       `json:"thisServed"`
	TotalServed int                       `json:"totalServed"`
}

// BunnyCommand will fetch a random picture from shibe.online
func BunnyCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		apiURL := "https://api.bunnies.io/v2/loop/random/?media=gif"
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

		var result BunniesIOResponse

		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &result)

		image := &tb.Document{File: tb.FromURL(result.Media.Gif)}
		bot.Send(m.Chat, image)
	}
}
