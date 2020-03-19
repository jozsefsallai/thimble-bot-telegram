package utils

import (
	"fmt"
	"io"

	"github.com/jozsefsallai/gophersauce"
	"github.com/jozsefsallai/thimble-bot-telegram/config"
)

// ImageSourceLookup will take a reader and try to detect the source of the image
// inside of the reader.
func ImageSourceLookup(image io.Reader) (string, error) {
	apiKey := config.GetConfig().Bot.SaucenaoAPIKey
	if len(apiKey) == 0 {
		return "(unknown)", nil
	}

	client, err := gophersauce.NewClient(&gophersauce.Settings{
		APIKey: apiKey,
	})

	if err != nil {
		return "", err
	}

	response, err := client.FromReader(image)
	if err != nil {
		return "", err
	}

	if response.Count() == 0 {
		return "(unknown)", nil
	}

	first := response.First()
	urls := first.Data.ExternalURLs

	if len(urls) == 0 {
		return "(unknown)", nil
	}

	return fmt.Sprintf("%s (automatically detected using SauceNAO, might be inaccurate)", urls[0]), nil
}
