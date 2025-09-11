package msgparser_test

import (
	msgparser "chat_server/MsgParser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserSimpleMessage(t *testing.T) {
	msg := []byte("msg:bartek:1")
	message, err := msgparser.ParseMsg(msg)
	assert.NoError(t, err)
	assert.Equal(t, "msg", message.Type)
	assert.Equal(t, "bartek", message.Value)
	assert.Equal(t, uint32(1), message.UserId)
}

func TestParserInvalidType(t *testing.T) {
	msg := []byte("test:test:1")
	_, err := msgparser.ParseMsg(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), msgparser.INVALID_TYPE)
}

func TestMessageToString(t *testing.T) {
	message, err := msgparser.NewMessage(msgparser.MSG, "hello world", 1, 1)
	assert.NoError(t, err)
	msg := message.ToBytes()
	assert.Equal(t, []byte("msg:hello world:1:1"), msg)
}

func TestParserInvalidMessage(t *testing.T) {
	msg := []byte("msgbartek")
	message, err := msgparser.ParseMsg(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid type: msgbartek")
	assert.Equal(t, "", message.Value)
	assert.Equal(t, "", message.Type)
}

func TestParserOnlyType(t *testing.T) {
	msg := []byte("list_rooms")
	message, err := msgparser.ParseMsg(msg)
	assert.NoError(t, err)
	assert.Equal(t, msgparser.LIST_ROOMS, message.Type)
	assert.Equal(t, "", message.Value)
	assert.Equal(t, uint32(0), message.UserId)
	assert.Equal(t, []byte("list_rooms"), message.ToBytes())
}

func TestParserTypeAndValue(t *testing.T) {
	msg := []byte("join:general")
	message, err := msgparser.ParseMsg(msg)
	assert.NoError(t, err)
	assert.Equal(t, msgparser.JOIN, message.Type)
	assert.Equal(t, "general", message.Value)
	assert.Equal(t, uint32(0), message.UserId)
	assert.Equal(t, []byte("join:general"), message.ToBytes())
}

func TestParserTypeValueAndId(t *testing.T) {
	msg := []byte("msg:hello:42")
	message, err := msgparser.ParseMsg(msg)
	assert.NoError(t, err)
	assert.Equal(t, msgparser.MSG, message.Type)
	assert.Equal(t, "hello", message.Value)
	assert.Equal(t, uint32(42), message.UserId)
	assert.Equal(t, []byte("msg:hello:42"), message.ToBytes())
}

func TestParserEmptyValueWithId(t *testing.T) {
	msg := []byte("id::42")
	message, err := msgparser.ParseMsg(msg)
	assert.NoError(t, err)
	assert.Equal(t, msgparser.ID_REQUEST, message.Type)
	assert.Equal(t, "", message.Value)
	assert.Equal(t, uint32(42), message.UserId)
	assert.Equal(t, []byte("id::42"), message.ToBytes())
}
