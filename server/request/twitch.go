package request

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	base_api_url   string = "https://api.twitch.tv/"
	base_oauth_url string = "https://id.twitch.tv/oauth2/"
)

//////////////////////
//  OAUTH RELATION  //
//////////////////////

func ValidateTwitchToken(authToken string) *TwitchOAuthValidation {
	return perform[TwitchOAuthValidation](true, Request{
		Method: "GET",
		URL:    oauth2("validate"),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("OAuth %s", authToken),
		},
	})
}

func RefreshCurrentTwitchToken() *TwitchOAuthRefresh {
	twitchProfile := Profiles.Twitch
	return perform[TwitchOAuthRefresh](true, Request{
		Method: "POST",
		URL:    oauth2("token"),
		Query: map[string]string{
			"client_id":     twitchProfile.ClientID,
			"client_secret": twitchProfile.ClientSecret,
			"grant_type":    "refresh_token",
			"refresh_token": url.QueryEscape(twitchProfile.RefreshToken),
		},
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	})
}

func RevokeTwitchToken(token string) {
	perform[any](false, Request{
		Method: "POST",
		URL:    oauth2("revoke"),
		Query: map[string]string{
			"client_id": Profiles.Twitch.ClientID,
			"token":     url.QueryEscape(token),
		},
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	})
}

//////////////////////
//  USER RELATION   //
//////////////////////

func TwitchUsersBy(usernames string) *TwitchUserList {
	requestProfile := Profiles.Twitch
	return perform[TwitchUserList](true, Request{
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
}

func TwitchUserBy(username string) *TwitchUser {
	return TwitchUsersBy(username).First()
}

func (r *TwitchUserList) First() *TwitchUser {
	users := r.Users
	if len(users) > 0 {
		return &users[0]
	}
	return nil
}

//////////////////////
// HELPER FUNCTIONS //
//////////////////////

func perform[T interface{}](shouldParse bool, request Request) *T {
	req := Build(request)
	res, err := client.Do(req)

	if err != nil || !shouldParse || (res.StatusCode != 200 && res.StatusCode != 202) {
		return nil
	}

	var data T
	json.NewDecoder(res.Body).Decode(&data)
	return &data
}

func helix(path string) string {
	return fmt.Sprintf("%shelix%s", base_api_url, path)
}

func oauth2(path string) string {
	return fmt.Sprintf("%s%s", base_oauth_url, path)
}
