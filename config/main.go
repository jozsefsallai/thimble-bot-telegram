package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// BotConfig contains properties related to the general bot application
type BotConfig struct {
	Token          string `json:"token"`
	SaucenaoAPIKey string `json:"saucenaoAPIKey"`
}

// PermissionsConfig contains bot permissions for certain commands and features
type PermissionsConfig struct {
	CanUploadMIA []int `json:"canUploadMIA"`
}

// Config contains the properties of a Thimble Bot configuration
type Config struct {
	Bot         BotConfig         `json:"bot"`
	Permissions PermissionsConfig `json:"permissions"`
}

// GetConfig returns the current configuration
func GetConfig() Config {
	configPath := "./config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file does not exist. Please create a config.json")
		os.Exit(1)
	}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatal("Failed to open the config file.", err)
	}

	defer f.Close()

	contents, _ := ioutil.ReadAll(f)

	var config Config

	err = json.Unmarshal(contents, &config)
	if err != nil {
		log.Fatal("Failed to parse the JSON config file. Maybe your config is invalid?", err)
	}

	return config
}
