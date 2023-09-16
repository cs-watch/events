package message

import (
	"sync"

	"github.com/cs-watch/events/errors"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -destination mock/acknowledger.go . Acknowledger
type Acknowledger interface {
	Ack(multiple bool) error
	Reject(requeue bool) error
	Nack(multiple bool, requeue bool) error
}

type Message struct {
	m            sync.Mutex
	done         bool
	acknowledger Acknowledger
	topic        string
	body         []byte
}

func NewMessage(acknowledger Acknowledger, topic string, body []byte) *Message {
	return &Message{
		m:            sync.Mutex{},
		done:         false,
		acknowledger: acknowledger,
		topic:        topic,
		body:         body,
	}
}

func (m *Message) Topic() string {
	return m.topic
}

func (m *Message) Payload() []byte {
	return m.body
}

func (m *Message) Ack() error {
	m.m.Lock()
	defer m.m.Unlock()

	if m.done {
		return errors.ErrMessageAlreadyAcknowledged
	}

	err := m.acknowledger.Ack(false)
	if err != nil {
		return err
	}
	m.done = true
	return nil
}

func (m *Message) Nack(requeue bool) error {
	m.m.Lock()
	defer m.m.Unlock()

	if m.done {
		return errors.ErrMessageAlreadyAcknowledged
	}

	err := m.acknowledger.Nack(false, requeue)
	if err != nil {
		return err
	}

	m.done = true
	return nil
}

func (m *Message) Reject(requeue bool) error {
	m.m.Lock()
	defer m.m.Unlock()

	if m.done {
		return errors.ErrMessageAlreadyAcknowledged
	}

	err := m.acknowledger.Reject(requeue)
	if err != nil {
		return err
	}

	m.done = true
	return nil
}
