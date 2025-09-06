package msgparser

import (
	"fmt"
	"strings"
)

const (
	NEW_CONNECTION = "new"
	RECONNECT      = "reconnect"
	ACCEPT         = "accept"
	MSG            = "msg"
	DISCONNECT     = "disconnect"
	ID_REQUEST     = "id"
	INVALID_MSG    = "invalid Message"
	INVALID_TYPE   = "invalid type"
)

type Message struct {
	Type  string
	Value string
}

func NewMessage(mType, mValue string) (Message, error) {
	if !ValidateType(mType) {
		return Message{Type: "", Value: ""}, fmt.Errorf("%s: %s", INVALID_TYPE, mType)
	}
	return Message{
		Type:  mType,
		Value: mValue,
	}, nil
}

func ParseMsg(msg []byte) (Message, error) {
	str := strings.SplitN(string(msg), ":", 2)
	if len(str) == 1 {
		return Message{}, fmt.Errorf("invalid message format: missing ':' separator in '%s'", string(msg))
	}
	Message, err := NewMessage(str[0], str[1])
	return Message, err
}

func ValidateType(str string) bool {
	switch str {
	case NEW_CONNECTION, RECONNECT, MSG, DISCONNECT, ID_REQUEST, ACCEPT:
		return true
	default:
		return false
	}
}

func (m Message) ToBytes() []byte {
	return []byte(m.Type + ":" + m.Value)
}
