package service

import "context"

type MessageProducer interface {
	SendMessage(ctx context.Context, value []byte) error
}
