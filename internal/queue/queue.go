package queue

import "context"

type Producer interface {
	Produce(ctx context.Context, message []byte) error
	Close() error
}

type Consumer interface {
	Consume() (<-chan string, error)
	Close() error
}
