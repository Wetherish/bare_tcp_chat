package msgparser_test

import (
	msgparser "chat_server/MsgParser"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserSimpleMessage(t *testing.T) {
	msg := []byte("msg:bartek")
	message, err := msgparser.ParseMsg(msg)
	assert.Equal(t, nil, err)
	assert.Equal(t, "msg", message.Type)
	assert.Equal(t, "bartek", message.Value)
}

func TestParserInvalidType(t *testing.T) {
	msg := []byte("test:test")
	_, err := msgparser.ParseMsg(msg)
	assert.Error(t, errors.New(msgparser.INVALID_TYPE), err)
}

func TestMessageToString(t *testing.T) {
	message, err := msgparser.NewMessage(msgparser.MSG, "hello world")
	msg := message.ToBytes()
	assert.Nil(t, err)
	assert.Equal(t, []byte("msg:hello world"), msg)
}

func TestParserInvalidMessage(t *testing.T) {
	msg := []byte("msgbartek")
	message, err := msgparser.ParseMsg(msg)
	assert.Error(t, errors.New(msgparser.INVALID_MSG), err)
	assert.Equal(t, "", message.Value)
	assert.Equal(t, "", message.Type)
}

