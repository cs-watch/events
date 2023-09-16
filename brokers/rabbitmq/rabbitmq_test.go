package rabbitmq

import (
	"context"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cs-watch/events/testhelpers"
	"github.com/cs-watch/events/testhelpers/containers"
)

func TestRabbitBroker_Subscribe(t *testing.T) {
	ctx := context.Background()

	container, url := containers.Prepare(ctx, t)
	defer container.Terminate(ctx)

	conn, err := amqp091.Dial(url)
	require.NoError(t, err)

	expectedTopic := "test"
	expectedBody := `{"test":"test"}`
	testhelpers.DeclareTopic(t, conn, expectedTopic)
	testhelpers.SendMessage(t, conn, expectedTopic, expectedBody)

	broker := New(conn)
	err = broker.Subscribe(ctx, expectedTopic)
	require.NoError(t, err)

	messageC := broker.Messages(ctx)
	msg := <-messageC
	assert.Equal(t, msg.Payload(), []byte(expectedBody))
}
