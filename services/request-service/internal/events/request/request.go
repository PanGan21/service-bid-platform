package request

import (
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/messaging"
)

type requestEvents struct {
	pub messaging.Publisher
}

const (
	REQUEST_CREATED_TOPIC = "request-created"
)

func NewRequestEvents(pub messaging.Publisher) *requestEvents {
	return &requestEvents{pub: pub}
}

func (events *requestEvents) PublishRequestCreated(request *entity.Request) error {
	msg := messaging.Message{
		Payload: request,
	}

	err := events.pub.Publish(REQUEST_CREATED_TOPIC, msg)
	if err != nil {
		return err
	}

	return nil
}
