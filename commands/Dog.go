package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tb "gopkg.in/tucnak/telebot.v2"
)

// RandomDogAPIResponse contains the JSON response returned by random.dog
type RandomDogAPIResponse struct {
	FileSizeBytes int `json:"fileSizeBytes"`
	URL string `json:"url"`
}

// DogCommand will fetch a random picture, GIF, or video from random.dog
func DogCommand(bot *tb.Bot) interface{} {
	return func(m *tb.Message) {
		apiURL := "https://random.dog/woof.json"
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

		var result RandomDogAPIResponse

		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &result)

		document := &tb.Document{File: tb.FromURL(result.URL)}
		bot.Send(m.Chat, document)
	}
}
