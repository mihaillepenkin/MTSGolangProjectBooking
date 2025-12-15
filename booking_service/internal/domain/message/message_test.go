package message

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage("123@mail.com", "1", "", "", StatusOK)

	assert.Equal(t, msg.Status, StatusOK)
}
