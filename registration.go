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

func (r *Application) HandleRegistration(
	response http.ResponseWriter,
	request *http.Request,
	_ httprouter.Params,
) {
	var requestBody TempRegistrationBody
	bodyErr := json.NewDecoder(request.Body).Decode(&requestBody)

	if bodyErr != nil {
		RespondJson(response, http.StatusInternalServerError, SomethingWentWrong, nil)
		return
	}

	email := strings.ToLower(requestBody.Email)
	username := strings.ToLower(requestBody.Username)

	// validate email
	if email == "" {
		RespondJson(response, http.StatusBadRequest, EmailMissing, nil)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		RespondJson(response, http.StatusBadRequest, EmailInvalid, nil)
		return
	}

	// validate requirements for username
	if username == "" {
		RespondJson(response, http.StatusBadRequest, UsernameMissing, nil)
		return
	}

	if usernameLen := len(username); usernameLen < 6 || usernameLen > 16 {
		RespondJson(response, http.StatusBadRequest, UsernameRequirementsNotMet, nil)
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
		RespondJson(response, http.StatusBadRequest, UsernameInvalid, nil)
		return
	}

	// ensure that username and email has not been taken already
	userLookup := new(User)
	r.DbClient.Model(userLookup).Where("email = ?", email).WhereOr("username = ?", username)

	if userLookup != nil {
		if userLookup.Name == username && userLookup.Email == email {
			RespondJson(response, http.StatusBadRequest, UsernameAndEmailAlreadyInUse, nil)
			return
		}

		if userLookup.Name == username {
			RespondJson(response, http.StatusBadRequest, UsernameAlreadyInUse, nil)
			return
		}

		if userLookup.Email == email {
			RespondJson(response, http.StatusBadRequest, EmailAlreadyInUse, nil)
			return
		}
	}

	// validate requirements for password
	password := requestBody.Password

	if password == "" {
		RespondJson(response, http.StatusBadRequest, PasswordMissing, nil)
		return
	}

	if passwordLen := len(password); passwordLen < 6 {
		RespondJson(response, http.StatusBadRequest, PasswordTooSmall, nil)
		return
	}
}
