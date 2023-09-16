package message

import (
	"testing"

	"github.com/cs-watch/events"
	"github.com/cs-watch/events/errors"
	"github.com/cs-watch/events/message/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMessage_Acknowledge(t *testing.T) {
	var (
		expectedTopic   = "topic"
		expectedPayload = []byte("payload")
	)
	testCases := []struct {
		name    string
		prepare func(m *mock_rabbitmq.MockAcknowledger)
		action  func(t *testing.T, msg events.Message)
	}{
		{
			name: "success positive confirmation",
			prepare: func(m *mock_rabbitmq.MockAcknowledger) {
				m.EXPECT().Ack(gomock.Any()).Return(nil).Times(1)
			},
			action: func(t *testing.T, msg events.Message) {
				err := msg.Ack()
				assert.NoError(t, err)
			},
		},
		{
			name: "success negative confirmation",
			prepare: func(m *mock_rabbitmq.MockAcknowledger) {
				m.EXPECT().Nack(false, false).Return(nil).Times(1)
			},
			action: func(t *testing.T, msg events.Message) {
				err := msg.Nack(false)
				assert.NoError(t, err)
			},
		},
		{
			name: "error multiple positive confirmations",
			prepare: func(m *mock_rabbitmq.MockAcknowledger) {
				m.EXPECT().Ack(false).Return(nil).Times(1)
			},
			action: func(t *testing.T, msg events.Message) {
				err := msg.Ack()
				assert.NoError(t, err)

				err = msg.Ack()
				assert.ErrorIs(t, err, errors.ErrMessageAlreadyAcknowledged)
			},
		},
		{
			name: "error multiple negative confirmations",
			prepare: func(m *mock_rabbitmq.MockAcknowledger) {
				m.EXPECT().Nack(false, false).Return(nil).Times(1)
			},
			action: func(t *testing.T, msg events.Message) {
				err := msg.Nack(false)
				assert.NoError(t, err)

				err = msg.Nack(false)
				assert.ErrorIs(t, err, errors.ErrMessageAlreadyAcknowledged)
			},
		},
		{
			name: "expected topic",
			prepare: func(m *mock_rabbitmq.MockAcknowledger) {
			},
			action: func(t *testing.T, msg events.Message) {
				assert.Equal(t, msg.Topic(), expectedTopic)
			},
		},
		{
			name: "expected payload",
			prepare: func(m *mock_rabbitmq.MockAcknowledger) {
			},
			action: func(t *testing.T, msg events.Message) {
				assert.Equal(t, msg.Payload(), expectedPayload)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mock_rabbitmq.NewMockAcknowledger(ctrl)
			tc.prepare(mock)
			msg := NewMessage(mock, expectedTopic, expectedPayload)
			tc.action(t, msg)
		})
	}
}
