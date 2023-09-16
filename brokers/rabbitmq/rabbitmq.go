package rabbitmq

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/cs-watch/events"
	"github.com/cs-watch/events/brokers"
	"github.com/cs-watch/events/message"
	"github.com/rabbitmq/amqp091-go"
)

type channel struct {
	channel *amqp091.Channel
	errorsC chan *amqp091.Error
}

type RabbitBroker struct {
	conn      *amqp091.Connection
	channels  map[string]channel
	messagesC chan events.Message
	errorsC   chan error
}

func New(conn *amqp091.Connection) *RabbitBroker {
	return &RabbitBroker{
		conn:      conn,
		channels:  make(map[string]channel),
		messagesC: make(chan events.Message, 1),
		errorsC:   make(chan error, 1),
	}
}

func (r *RabbitBroker) Subscribe(ctx context.Context, topic string) error {
	c, err := r.conn.Channel()
	if err != nil {
		return errors.Join(err, brokers.ErrFailedToOpenChannel)
	}

	go func(ctx context.Context, c *amqp091.Channel) {
		<-ctx.Done()
		if err := c.Close(); err != nil {
			r.errorsC <- errors.Join(err, brokers.ErrChannelIsClosed)
		}
	}(ctx, c)

	r.channels[topic] = channel{
		channel: c,
		errorsC: make(chan *amqp091.Error, 1),
	}

	go func() {
		log.Printf("Closing: %s", <-c.NotifyClose(r.channels[topic].errorsC))
	}()

	go func(ctx context.Context, c channel) {
		delivery, err := c.channel.Consume(
			topic,
			"",
			false, // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,
		)
		if err != nil {
			slog.Error("failed to consume", "error", err)
			r.errorsC <- err
		}

		for {
			select {
			case <-ctx.Done():
				close(r.messagesC)
				return
			case rawMsg, channelAlive := <-delivery:
				if !channelAlive {
					slog.Warn("rabbitmq.Listen channel is closed!")
					return
				}
				msg := message.NewMessage(rawMsg, topic, rawMsg.Body)
				r.messagesC <- msg
			}
		}
	}(ctx, r.channels[topic])

	return nil
}

func (r *RabbitBroker) Unsubscribe(ctx context.Context, topic string) error {
	if err := r.channels[topic].channel.Close(); err != nil {
		return err
	}
	delete(r.channels, topic)
	return nil
}

func (r *RabbitBroker) Messages(ctx context.Context) <-chan events.Message {
	return r.messagesC
}

func (r *RabbitBroker) Errors(ctx context.Context) <-chan error {
	errC := make(chan error, len(r.channels))
	for _, c := range r.channels {
		go func(c channel) {
			e, open := <-c.errorsC
			if !open {
				slog.Warn("rabbitmq.Errors channel with errors is closed!")
			}

			errC <- errors.New(e.Error())
		}(c)
	}
	return r.errorsC
}
