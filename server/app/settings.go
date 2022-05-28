package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/imoliwer/sound-point-twitch-bot/server/request"
)

type TwitchCommandOption struct {
	Enabled bool `json:"enabled"`
}

type TwitchCommandPrimaryOption struct {
	TwitchCommandOption
	Arguments map[string]TwitchCommandOption `json:"arguments"`
}

type TwitchCommandMessages struct {
	PointsNoArg            string `json:"points_no_arg"`
	PointsGiveSuccess      string `json:"points_give_success"`
	PointsSetSuccess       string `json:"points_set_success"`
	SpecifyUser            string `json:"specify_user"`
	SpecifyAmount          string `json:"specify_amount"`
	MustSpecifyValidAmount string `json:"must_specify_valid_amount"`
	CouldNotFindUser       string `json:"could_not_find_user"`
}

type TwitchCommandSettings struct {
	Prefix   string                                `json:"prefix"`
	Options  map[string]TwitchCommandPrimaryOption `json:"options"`
	Messages TwitchCommandMessages                 `json:"messages"`
}

type TwitchBotSettings struct {
	Name      string                `json:"name"`
	AuthToken string                `json:"auth_token"`
	Channel   string                `json:"channel_to_join"`
	Command   TwitchCommandSettings `json:"command"`
}

type TempTwitchAccessSettings struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

type Settings struct {
	TwitchBot       TwitchBotSettings         `json:"twitch_chat_bot"`
	TwitchAccessory *TempTwitchAccessSettings `json:"twitch_accessories"` // temporary
}

func (r *Settings) Save() {
	profile := request.Profiles.Twitch
	r.TwitchAccessory = &TempTwitchAccessSettings{
		ClientID:     profile.ClientID,
		ClientSecret: profile.ClientSecret,
		AuthToken:    profile.OAuthToken,
		RefreshToken: profile.RefreshToken,
	}

	bytes, _ := json.MarshalIndent(r, "", "  ")
	os.WriteFile("settings.json", bytes, 0)
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
			"prefix": "!",
			"options": {
				"points": {
					"enabled": true,
					"arguments": {
						"set": {
							"enabled": true
						},
						"give": {
							"enabled": true
						}
					}
				}
			},
			"messages": {
				"points_no_arg": "You currently have %d points.",
				"points_give_success": "%s has been given %d points.",
				"points_set_success": "The points of %s has been set to %d.",
				"specify_user": "You must specify a user.",
				"specify_amount": "You must specify an amount.",
				"must_specify_valid_amount": "You must specify a valid amount.",
				"could_not_find_user": "Could not find the user %s."
			}
		}
	},
	"twitch_accessories": {
		"client_id": "<your_client_id>",
		"client_secret": "<your_client_secret>",
		"auth_token": "<your_auth_token>",
		"refresh_token": "<your_refresh_token>",
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
