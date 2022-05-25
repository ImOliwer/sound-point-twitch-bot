package request

import (
	"encoding/json"
	"fmt"
)

const (
	base_url string = "https://api.twitch.tv/"
)

func TwitchUsersBy(usernames string) *TwitchUserList {
	requestProfile := request_profiles.Twitch
	request := Build(Request{
		Method: "GET",
		URL:    helix("/users"),
		Query: map[string]string{
			"login": usernames,
		},
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", requestProfile.OAuthToken),
			"Client-ID":     requestProfile.ClientID,
		},
	})

	response, err := client.Do(request)
	if err != nil {
		return nil
	}

	var list TwitchUserList
	json.NewDecoder(response.Body).Decode(&list)
	return &list
}

func TwitchUserBy(username string) *TwitchUser {
	return TwitchUsersBy(username).First()
}

func helix(path string) string {
	return fmt.Sprintf("%shelix%s", base_url, path)
}

func (r *TwitchUserList) First() *TwitchUser {
	users := r.Users
	if len(users) > 0 {
		return &users[0]
	}
	return nil
}
