package msgparser

import (
	"bytes"
	"errors"
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
	UserId  uint32
	Type    string
	Content []byte
	RoomId  uint32
}

func NewMessage(msg_type string, content []byte, id, roomId uint32) (Message, error) {
	if !ValidateType(msg_type) {
		return Message{}, fmt.Errorf("%s: %s", INVALID_TYPE, msg_type)
	}
	return Message{
		Type:    msg_type,
		UserId:  id,
		RoomId:  roomId,
		Content: content,
	}, nil
}

func ParseMsg(msg []byte) (Message, error) {
	string_msg := string(msg)
	elems := strings.SplitN(string_msg, ":", 4)

	if len(elems) < 4 {
		return Message{}, errors.New("invalid message format")
	}

	msg_type := elems[0]
	if !ValidateType(msg_type) {
		return Message{}, errors.New(INVALID_TYPE)
	}

	id, err := strconv.Atoi(elems[1])
	if err != nil {
		return Message{}, errors.New("invalid Id")
	}
	roomID, err := strconv.Atoi(elems[2])
	if err != nil {
		return Message{}, errors.New("invalid Id")
	}

	content := []byte(elems[3])

	return Message{UserId: uint32(id), Type: msg_type, RoomId: uint32(roomID), Content: content}, nil

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
	var buffor bytes.Buffer
	buffor.WriteString(fmt.Sprintf("%s:%d:%d:", m.Type, m.UserId, m.RoomId))
	buffor.Write(m.Content)
	return buffor.Bytes()
}
