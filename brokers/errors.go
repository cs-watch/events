package brokers

import "errors"

var (
	ErrFailedToOpenChannel = errors.New("failed to open channel")
	ErrChannelIsClosed     = errors.New("channel is closed")
)
