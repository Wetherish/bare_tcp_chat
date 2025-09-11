package msgparser

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ACCEPT       = "accept"
	MSG          = "msg"
	JOIN         = "join"
	LEAVE        = "leave"
	LIST_ROOMS   = "list_rooms"
	ID_REQUEST   = "id"
	INVALID_MSG  = "invalid message"
	INVALID_TYPE = "invalid type"
)

type Message struct {
	UserId uint32
	Type   string
	Value  string
	RoomId uint32
}

func NewMessage(mType, mValue string, id, roomId uint32) (Message, error) {
	if !ValidateType(mType) {
		return Message{}, fmt.Errorf("%s: %s", INVALID_TYPE, mType)
	}
	return Message{
		UserId: id,
		Type:   mType,
		Value:  mValue,
		RoomId: roomId,
	}, nil
}

func ParseMsg(msg []byte) (Message, error) {
	parts := strings.SplitN(string(msg), ":", 4)

	var mType, mValue string
	var id, roomId uint32
	var err error

	switch len(parts) {
	case 1:
		mType = parts[0]
	case 2:
		mType = parts[0]
		mValue = parts[1]
	case 3:
		mType = parts[0]
		mValue = parts[1]
		if parts[2] != "" {
			var parsed int64
			parsed, err = strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return Message{}, fmt.Errorf("invalid user id: %v", err)
			}
			id = uint32(parsed)
		}
	case 4:
		mType = parts[0]
		mValue = parts[1]
		if parts[2] != "" {
			var parsed int64
			parsed, err = strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return Message{}, fmt.Errorf("invalid user id: %v", err)
			}
			id = uint32(parsed)
		}
		if parts[3] != "" {
			var parsedRoom int64
			parsedRoom, err = strconv.ParseInt(parts[3], 10, 32)
			if err != nil {
				return Message{}, fmt.Errorf("invalid room id: %v", err)
			}
			roomId = uint32(parsedRoom)
		}
	default:
		return Message{}, fmt.Errorf("invalid message format: %s", string(msg))
	}

	return NewMessage(mType, mValue, id, roomId)
}

func ValidateType(str string) bool {
	switch str {
	case ID_REQUEST, MSG, ACCEPT, JOIN, LEAVE, LIST_ROOMS:
		return true
	default:
		return false
	}
}

func (m Message) ToBytes() []byte {
	switch {
	case m.Value == "" && m.UserId == 0 && m.RoomId == 0:
		return []byte(m.Type)
	case m.UserId == 0 && m.RoomId == 0:
		return []byte(fmt.Sprintf("%s:%s", m.Type, m.Value))
	case m.RoomId == 0:
		return []byte(fmt.Sprintf("%s:%s:%d", m.Type, m.Value, m.UserId))
	default:
		return []byte(fmt.Sprintf("%s:%s:%d:%d", m.Type, m.Value, m.UserId, m.RoomId))
	}
}
