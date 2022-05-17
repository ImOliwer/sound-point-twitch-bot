package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type BotSettings struct {
	Name      string `json:"name"`
	AuthToken string `json:"auth_token"`
}

type Settings struct {
	Bot BotSettings `json:"bot"`
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
  "bot": {
		"name": "bot_name_here",
		"auth_token": "bot_auth_token_here"
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
