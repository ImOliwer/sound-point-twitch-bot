package main

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func RespondJson(writer http.ResponseWriter, status int, message string, data any) error {
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(JsonResponse{Message: message, Data: data})
}
