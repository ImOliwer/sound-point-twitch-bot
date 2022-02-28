package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type TempRegistrationBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleRegistration(
	response http.ResponseWriter,
	request *http.Request,
	_ httprouter.Params,
) {
	var requestBody TempRegistrationBody
	bodyErr := json.NewDecoder(request.Body).Decode(&requestBody)

	if bodyErr != nil {
		RespondJson(response, http.StatusInternalServerError, "something went wrong...", nil)
		return
	}

	email := strings.ToLower(requestBody.Email)
	username := strings.ToLower(requestBody.Username)

	// validate email
	if email == "" {
		RespondJson(response, http.StatusBadRequest, "email is missing", nil)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		RespondJson(response, http.StatusBadRequest, "email is invalid", nil)
		return
	}

	// validate requirements for username
	if username == "" {
		RespondJson(response, http.StatusBadRequest, "username is missing", nil)
		return
	}

	if usernameLen := len(username); usernameLen < 6 || usernameLen > 16 {
		RespondJson(response, http.StatusBadRequest, "username must contain 6-16 characters", nil)
		return
	}

	usernameAlphaFailure := false
	for _, letter := range username {
		if letter >= 'a' && letter <= 'z' || letter >= 'A' && letter <= 'Z' || letter >= '0' && letter <= '9' || letter == '_' {
			continue
		}
		usernameAlphaFailure = true
		break
	}

	if usernameAlphaFailure {
		RespondJson(response, http.StatusBadRequest, "username is invalid", nil)
		return
	}

	// TODO: ensure that username and email has not been taken already

	// validate requirements for password
	password := requestBody.Password

	if password == "" {
		RespondJson(response, http.StatusBadRequest, "password is missing", nil)
		return
	}

	if passwordLen := len(password); passwordLen < 6 {
		RespondJson(response, http.StatusBadRequest, "password too small", nil)
		return
	}
}
