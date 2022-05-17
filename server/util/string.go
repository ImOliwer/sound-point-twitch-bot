package util

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type FormatableString struct {
	Data      string
	Arguments []any
}

func NewFormatableString(data string, args ...any) FormatableString {
	return FormatableString{Data: data, Arguments: args}
}

func SendMultipleString(connection *websocket.Conn, formatables []FormatableString) {
	for _, formatable := range formatables {
		arguments := formatable.Arguments
		if arguments == nil {
			SendString(connection, formatable.Data)
			continue
		}
		SendString(connection, formatable.Data, arguments...)
	}
}

func SendString(connection *websocket.Conn, data string, args ...any) {
	websocket.Message.Send(connection, fmt.Sprintf(data, args...))
}
