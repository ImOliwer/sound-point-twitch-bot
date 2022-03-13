package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LoginProps struct {
	Token       string `json:"token"`
	Credentials *struct {
		Name     string
		Password string
	} `json:"credentials"`
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
}
