package request

import (
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/messaging"
)

type requestEvents struct {
	pub messaging.Publisher
}

func NewRequestEvents(pub messaging.Publisher) *requestEvents {
	return &requestEvents{pub: pub}
}

func (events *requestEvents) PublishRequestCreated(topic string, request *entity.Request) error {
	msg := messaging.Message{
		Payload: request,
	}

	err := events.pub.Publish(topic, msg)
	if err != nil {
		return err
	}

	return nil
}
