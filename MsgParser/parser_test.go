package msgparser_test

import (
	msgparser "chat_server/MsgParser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserSimpleMessage(t *testing.T) {
	msg := []byte("msg:1:0:bartek")
	message, err := msgparser.ParseMsg(msg)
	assert.NoError(t, err)
	assert.Equal(t, "msg", message.Type)
	assert.Equal(t, []byte("bartek"), message.Content)
	assert.Equal(t, uint32(1), message.UserId)
}

func TestParserInvalidType(t *testing.T) {
	msg := []byte("test:1:1:test")
	_, err := msgparser.ParseMsg(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), msgparser.INVALID_TYPE)
}

func TestMessageToString(t *testing.T) {
	message, err := msgparser.NewMessage(msgparser.MSG, []byte("hello world"), 1, 1)
	assert.NoError(t, err)
	msg := message.ToBytes()
	assert.Equal(t, []byte("msg:1:1:hello world"), msg)
}

func TestParserInvalidMessage(t *testing.T) {
	msg := []byte("msgbartek")
	message, err := msgparser.ParseMsg(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid message format")
	assert.Equal(t, []byte(nil), message.Content)
	assert.Equal(t, "", message.Type)
}
