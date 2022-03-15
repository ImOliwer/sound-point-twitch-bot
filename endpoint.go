package main

import (
	"encoding/json"
	"net/http"
)

type ResponseErrorCode = uint8

const (
	// general
	OK ResponseErrorCode = iota
	SomethingWentWrong

	// bundle
	UsernameAndEmailAlreadyInUse
	UserFailedToCreate
	MissingCredentialsProperty
	EmailAndPasswordRequired

	// email
	EmailMissing
	EmailInvalid
	EmailAlreadyInUse

	// username
	UsernameMissing
	UsernameInvalid
	UsernameRequirementsNotMet
	UsernameAlreadyInUse

	// password
	PasswordMissing
	PasswordTooSmall
	PasswordFailedHash

	// login
	EmailAndPasswordMismatch
)

type JsonResponse struct {
	Data      any               `json:"data"`
	ErrorCode ResponseErrorCode `json:"errorCode"`
}

func RespondJson(writer http.ResponseWriter, status int, errorCode ResponseErrorCode, data any) error {
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(JsonResponse{Data: data, ErrorCode: errorCode})
}
