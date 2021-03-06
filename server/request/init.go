package request

import (
	"io"
	"net/http"
	"time"
)

type TwitchRequestProfile struct {
	ClientID     string
	ClientSecret string
	OAuthToken   string
	RefreshToken string
}

type RequestProfiles struct {
	Twitch TwitchRequestProfile
}

type Request struct {
	Method  string
	URL     string
	Body    io.Reader
	Query   map[string]string
	Headers map[string]string
}

var Profiles *RequestProfiles
var client = &http.Client{Timeout: time.Second * 5}

func Build(r Request) *http.Request {
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		panic(err)
	}

	query := req.URL.Query()
	for_build(
		r.Query,
		func(key string, value string) { query.Add(key, value) },
		func() { req.URL.RawQuery = query.Encode() },
	)

	for_build(r.Headers, func(key string, value string) { req.Header.Set(key, value) }, nil)
	return req
}

func for_build(with map[string]string, handle func(string, string), completion func()) {
	if with == nil || len(with) == 0 {
		return
	}

	for key, value := range with {
		handle(key, value)
	}

	if completion != nil {
		completion()
	}
}
