package events

import (
	"context"
)

type Message interface {
	Topic() string
	Payload() []byte

	Ack() error
	Nack(requeue bool) error
	Reject(requeue bool) error
}

type Broker interface {
	Subscribe(ctx context.Context, topic string) error
	Unsubscribe(ctx context.Context, topic string) error

	Messages(ctx context.Context) <-chan Message
	Errors(ctx context.Context) <-chan error
}

type Router interface {
	Subscribe(ctx context.Context, topic string, handler func(ctx context.Context, msg Message) error) error // Overwrites the previous handler for the topic
	Unsubscribe(ctx context.Context, topic string) error
	Listen(ctx context.Context, msgs <-chan Message) error
}

type Sender interface {
	Send(ctx context.Context, topic string, payload []byte) error
}
