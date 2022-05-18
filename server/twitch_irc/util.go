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
	indirectType := indirect.Type()
	var data map[string]string

	rawJson, _ := twitch_msg_to_json(twitchIrcMessage)
	json.Unmarshal([]byte(rawJson), &data)

	for index := 0; index < indirectType.NumField(); index++ {
		field := indirectType.Field(index)
		linkedTo, isLinked := field.Tag.Lookup("link")
		valueField := indirect.Field(index)

		if !isLinked {
			tag := field.Tag.Get("json")
			if tag == "" {
				continue
			}
			value := value_with_handler(tag, data[tag], handlers)
			valueField.Set(reflect.ValueOf(value))
			continue
		}

		newValue := reflect.New(field.Type).Elem()
		for _, rawLink := range strings.Split(linkedTo, ";") {
			linkSplit := strings.Split(rawLink, "=")
			linkJsonTag := linkSplit[0]
			linkField := linkSplit[1]
			value := value_with_handler(linkJsonTag, data[linkJsonTag], handlers)
			field := newValue.FieldByName(linkField)
			field.Set(reflect.ValueOf(value))
		}

		valueField.Set(newValue)
	}
	return nil
}

func value_with_handler(tag string, data string, handlers map[string]objectify_handler) interface{} {
	if handler, ok := handlers[tag]; ok {
		return handler(data)
	}
	return data
}
