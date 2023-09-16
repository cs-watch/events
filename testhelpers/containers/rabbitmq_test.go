package containers

import (
	"context"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestPrepare(t *testing.T) {
	ctx := context.Background()
	container, url := Prepare(ctx, t)
	defer container.Terminate(ctx)

	conn, err := amqp091.Dial(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	assert.False(t, conn.IsClosed())
}
