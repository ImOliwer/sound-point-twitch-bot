package twitch_irc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type objectify_handler = func(string) interface{}

func twitch_msg_to_json(twitchMsg string) (string, error) {
	length := len(twitchMsg)
	if length <= 1 {
		return "", errors.New("invalid twitch message")
	}

	if twitchMsg[0] == '@' {
		twitchMsg = twitchMsg[1:length]
	}

	builder := strings.Builder{}
	builder.WriteRune('{')

	isFirst := true
	for _, entry := range strings.Split(twitchMsg, ";") {
		splitEntry := strings.Split(entry, "=")
		key := splitEntry[0]
		rawValue := splitEntry[1]
		if isFirst {
			isFirst = false
		} else {
			builder.WriteRune(',')
		}
		builder.WriteString(fmt.Sprintf("\"%[1]s\": \"%[2]s\"", key, rawValue))
	}

	builder.WriteRune('}')
	return builder.String(), nil
}

// `json:"linked_item" link:"prop_one=Prop"` <- example tag
func objectify_irc(twitchIrcMessage string, toPtr interface{}, handlers map[string]objectify_handler) error {
	if toPtr == nil {
		return errors.New("to-pointer must not be nil")
	}

	indirect := reflect.Indirect(reflect.ValueOf(toPtr))
	var data map[string]string

	rawJson, _ := twitch_msg_to_json(twitchIrcMessage)
	err := json.Unmarshal([]byte(strings.ReplaceAll(rawJson, "\\s", " ")), &data)

	if err != nil {
		panic(err)
	}

	handle_prop(data, handlers, indirect)
	return nil
}

func handle_prop(data map[string]string, handlers map[string]objectify_handler, value reflect.Value) {
	baseType := value.Type()
	for index := 0; index < value.NumField(); index++ {
		valueField := value.Field(index)
		structField := baseType.Field(index)

		if it, ok := structField.Tag.Lookup("twitchObj"); !ok || it != "true" {
			tag := structField.Tag.Get("json")
			if tag == "" {
				continue
			}

			rawData, ok := data[tag]
			if !ok || rawData == "" {
				continue
			}

			value := value_with_handler(tag, rawData, handlers)
			valueField.Set(reflect.ValueOf(value))
			continue
		}

		newValue := reflect.New(structField.Type).Elem()
		handle_prop(data, handlers, newValue)
		valueField.Set(newValue)
	}
}

func value_with_handler(tag string, data string, handlers map[string]objectify_handler) interface{} {
	if handler, ok := handlers[tag]; ok {
		return handler(data)
	}
	return data
}
