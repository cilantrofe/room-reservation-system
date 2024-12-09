package interfaces

import "context"

//go:generate mockgen -source=producer.go -destination=mocks/producer_mock.go -package=mocks
type MessageProducer interface {
	SendMessage(ctx context.Context, value []byte) error
}
