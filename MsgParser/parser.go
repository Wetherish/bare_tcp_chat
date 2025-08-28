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
	INVALID_MSG    = "invalid Message"
	INVALID_TYPE   = "invalid type"
)

type Message struct {
	Type  string
	Value string
}

func NewMessage(mType, mValue string) (Message, error) {
	if !ValidateType(mType) {
		return Message{Type: "", Value: ""}, errors.New(INVALID_TYPE)
	}
	return Message{
		Type:  mType,
		Value: mValue,
	}, nil
}

func ParseMsg(msg []byte) (Message, error) {
	str := strings.SplitN(string(msg), ":", 2)
	if len(str) == 1 {
		empytMsg, _ := NewMessage("", "")
		return empytMsg, errors.New(INVALID_MSG)
	}
	Message, err := NewMessage(str[0], str[1])
	return Message, err
}

func ValidateType(str string) bool {
	switch str {
	case NEW_CONNECTION, RECONNECT, MSG, DISCONNECT:
		return true
	default:
		return false
	}
}

func (m *Message) ToBytes() []byte {
	return []byte(m.Type + ":" + m.Value)
}
