package router

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cs-watch/events"
)

type Router struct {
	subscribes map[string]func(ctx context.Context, msg events.Message) error
	broker     events.Broker
}

func New(broker events.Broker) *Router {
	return &Router{
		subscribes: make(map[string]func(ctx context.Context, msg events.Message) error),
		broker:     broker,
	}
}

func (r *Router) Subscribe(ctx context.Context, topic string, handler func(ctx context.Context, msg events.Message) error) error {
	r.broker.Subscribe(ctx, topic)

	r.subscribes[topic] = handler
	return nil
}

func (r *Router) Unsubscribe(ctx context.Context, topic string) error {
	r.broker.Unsubscribe(ctx, topic)

	delete(r.subscribes, topic)
	return nil
}

func (r *Router) Listen(ctx context.Context, msgs <-chan events.Message) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, channelAlive := <-msgs:
			if !channelAlive {
				slog.Warn("router.Listen channel is closed!")
				return nil
			}
			slog.Debug(fmt.Sprintf("%s: %s", msg.Topic(), msg.Payload()))
			handler := r.subscribes[msg.Topic()]
			err := handler(ctx, msg)
			if err != nil {
				return msg.Nack(false)
			}
			msg.Ack()
		}
	}
}
