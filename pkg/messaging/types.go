package messaging

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Options struct {
	Topic   string
	GroupId string
}

type Message struct {
	Id uuid.UUID
}

type BrokerMessage interface {
	ack() error
	nack() error
}

type EventHandler func(msg BrokerMessage) error

type Consumer interface {
	Close() error
}

type MessagingService interface {
	Subscribe(ctx context.Context, string, groupId string, handler EventHandler) *Consumer
	Publish(ctx context.Context, topic string, msg Message) error
}
