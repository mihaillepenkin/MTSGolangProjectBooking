package message

import "context"

type Producer interface {
	SendMessage(ctx context.Context, message *Message) error
}
