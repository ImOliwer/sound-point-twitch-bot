package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type LoginProps struct {
	Credentials *struct {
		Email    string
		Password string
	} `json:"credentials"`
}

type LoginSuccess struct {
	Token string `json:"token"`
}

func (r *Application) HandleLogin(
	response http.ResponseWriter,
	request *http.Request,
	_ httprouter.Params,
) {
	var requestBody LoginProps
	bodyErr := json.NewDecoder(request.Body).Decode(&requestBody)

	if bodyErr != nil {
		RespondJson(response, http.StatusInternalServerError, SomethingWentWrong, nil)
		return
	}

	credentials := requestBody.Credentials
	if credentials == nil {
		RespondJson(response, http.StatusBadRequest, MissingCredentialsProperty, nil)
		return
	}

	email := strings.ToLower(credentials.Email)
	password := credentials.Password

	if email == "" || password == "" {
		RespondJson(response, http.StatusBadRequest, EmailAndPasswordRequired, nil)
		return
	}

	found := new(User)
	_, lookupErr := r.DbClient.NewSelect().Model(found).Where("email = ?", email).Exec(context.Background())

	if lookupErr != nil {
		RespondJson(response, http.StatusBadRequest, EmailAndPasswordMismatch, nil)
		return
	}

	if !ComparePassword(found.Password, password) {
		RespondJson(response, http.StatusBadRequest, EmailAndPasswordMismatch, nil)
		return
	}

	token, tokenErr := r.GenerateJWT(found)
	if tokenErr != nil {
		RespondJson(response, http.StatusInternalServerError, SomethingWentWrong, nil)
		return
	}

	RespondJson(response, http.StatusOK, OK, LoginSuccess{Token: token})
}
