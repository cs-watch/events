package tests

import (
	"context"
	"log/slog"
	"sync"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cs-watch/events"
	"github.com/cs-watch/events/brokers/rabbitmq"
	"github.com/cs-watch/events/router"
	"github.com/cs-watch/events/testhelpers"
	"github.com/cs-watch/events/testhelpers/containers"
)

func TestEvents(t *testing.T) {
	ctx := context.Background()

	container, url := containers.Prepare(ctx, t)
	defer container.Terminate(ctx)

	conn, err := amqp091.Dial(url)
	require.NoError(t, err)

	expectedTopic := "test"
	expectedBody := `{"test":"test"}`
	testhelpers.DeclareTopic(t, conn, expectedTopic)

	broker := rabbitmq.New(conn)
	router := router.New(broker)

	var wg sync.WaitGroup

	wg.Add(1)
	router.Subscribe(ctx, expectedTopic, func(ctx context.Context, msg events.Message) error {
		slog.Debug("router.Subscribe(msg.payload)", "payload", msg.Payload())
		assert.Equal(t, expectedBody, string(msg.Payload()))
		wg.Done()
		return nil
	})

	go func() {
		err := router.Listen(ctx, broker.Messages(ctx))
		require.NoError(t, err)
	}()

	t.Run("handle message", func(t *testing.T) {
		testhelpers.SendMessage(t, conn, expectedTopic, expectedBody)
		require.NoError(t, err)
		wg.Wait()
	})
}
