package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type TwitchCommandSettings struct {
	Prefix string `json:"prefix"`
}

type TwitchBotSettings struct {
	Name      string                `json:"name"`
	AuthToken string                `json:"auth_token"`
	Channel   string                `json:"channel_to_join"`
	Command   TwitchCommandSettings `json:"command"`
}

type TempTwitchAccessSettings struct {
	ClientID     string `json:"client_id"`
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

type Settings struct {
	TwitchBot       TwitchBotSettings         `json:"twitch_chat_bot"`
	TwitchAccessory *TempTwitchAccessSettings `json:"twitch_accessories"` // temporary
}

func ReadSettings() *Settings {
	// bytes read from the settings file
	var settingsContent []byte

	// create file if it doesn't exist
	if _, settingsError := os.Stat("settings.json"); errors.Is(settingsError, os.ErrNotExist) {
		created, createdError := os.Create("settings.json")
		if createdError != nil {
			log.Panic("An error occurred during creation of 'settings.json' file. Try creating it manually.")
			return nil
		}
		settingsContent = []byte(`{
  "twitch_chat_bot": {
		"name": "<bot_username>",
		"auth_token": "<bot_auth_token>",
		"channel_to_join": "<your_channel_name>",
		"command": {
			"prefix": "!"
		}
	},
	"twitch_accessories": {
		"client_id": "<your_client_id>",
		"auth_token": "<your_auth_token>",
		"refresh_token": "<your_refresh_token>"
	}
}`)
		created.Write(settingsContent)
		created.Close()
		log.Println("No 'settings' file found. One was created for you, please modify it accordingly. Exiting...")
		time.Sleep(time.Second * 3)
		return nil
	}

	// fetch settings from file
	contentRead, _ := ioutil.ReadFile("settings.json")
	settingsContent = contentRead

	// unmarshal json into our var with corresponding struct
	var settings Settings
	json.Unmarshal(settingsContent, &settings)

	// return the settings accordingly
	return &settings
}
