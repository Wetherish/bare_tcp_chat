package msgparser

import (
	"fmt"
	"strings"
)

const (
	ACCEPT       = "accept"
	MSG          = "msg"
	ID_REQUEST   = "id"
	INVALID_MSG  = "invalid Message"
	INVALID_TYPE = "invalid type"
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
	case ID_REQUEST, MSG, ACCEPT:
		return true
	default:
		return false
	}
}

func (m Message) ToBytes() []byte {
	return []byte(m.Type + ":" + m.Value)
}
