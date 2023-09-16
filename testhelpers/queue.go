package testhelpers

import (
	"context"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func DeclareTopic(t *testing.T, conn *amqp091.Connection, topic string) {
	t.Helper()

	c, err := conn.Channel()
	require.NoError(t, err)

	queue, err := c.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	require.NoError(t, err)
	require.NotNil(t, queue)
}

func SendMessage(t *testing.T, conn *amqp091.Connection, topic string, body string) {
	t.Helper()

	c, err := conn.Channel()
	require.NoError(t, err)

	err = c.PublishWithContext(
		context.Background(),
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		},
	)
	require.NoError(t, err)
}
