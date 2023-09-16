package containers

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	Image           = "rabbitmq:3.12.4"
	DefaultUser     = "guest"
	DefaultPassword = "guest"
)

func Prepare(ctx context.Context, t *testing.T) (container testcontainers.Container, url string) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        Image,
		ExposedPorts: []string{"5672/tcp"},
		// WaitingFor:   wait.ForLog("Server startup complete;"),
		WaitingFor: wait.ForAll(
			wait.ForLog("Server startup complete;"),
			wait.ForListeningPort("5672/tcp"),
		),
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": DefaultUser,
			"RABBITMQ_DEFAULT_PASS": DefaultPassword,
		},
	}

	rC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Failed to create RabbitMQ container")

	mappedPort, err := rC.MappedPort(ctx, "5672/tcp")
	require.NoError(t, err, "Failed to get mapped port")

	ip, err := rC.Host(ctx)
	require.NoError(t, err, "Failed to get container IP")

	return rC, fmt.Sprintf("amqp://%s:%s@%s:%s/", DefaultUser, DefaultPassword, ip, mappedPort.Port())
}
