package main

import (
	"encoding/json"
	"net/http"
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
	
	email := strings.ToLower(requestBody.Email)
	username := strings.ToLower(requestBody.Username)

	if email == "" {
		RespondJson(response, http.StatusBadRequest, "email is missing", nil)
		return
	}

	// validate email

	if username == "" {
		RespondJson(response, http.StatusBadRequest, "username is missing", nil)
	}

	// validate requirements for username
}
