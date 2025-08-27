package msgparser

import (
	"errors"
	"strings"
)

const (
	NEW_CONNECTION = "new"
	RECONNECT      = "reconnect"
	ACCEPT         = "accept"
	MSG            = "msg"
	DISCONNECT     = "disconnect"
	INVALID_MSG    = "invalid message"
	INVALID_TYPE   = "invalid type"
)

type message struct {
	Type  string
	Value string
}

func NewMessage(mType, mValue string) (message, error) {
	if !ValidateType(mType) {
		return message{Type: "", Value: ""}, errors.New(INVALID_TYPE)
	}
	return message{
		Type:  mType,
		Value: mValue,
	}, nil
}

func ParseMsg(msg []byte) (message, error) {
	str := strings.SplitN(string(msg), ":", 2)
	if len(str) == 1 {
		empytMsg, _ := NewMessage("", "")
		return empytMsg, errors.New(INVALID_MSG)
	}
	message, err := NewMessage(str[0], str[1])
	return message, err
}

func ValidateType(str string) bool {
	switch str {
	case NEW_CONNECTION, RECONNECT, MSG, DISCONNECT:
		return true
	default:
		return false
	}
}

func (m *message) ToBytes() []byte {
	return []byte(m.Type + ":" + m.Value)
}
