package request

import (
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/messaging"
)

type requestEvents struct {
	pub messaging.Publisher
}

const (
	REQUEST_APPROVED_TOPIC = "request-approved"
	REQUEST_UPDATED_TOPIC  = "request-updated"
)

func NewRequestEvents(pub messaging.Publisher) *requestEvents {
	return &requestEvents{pub: pub}
}

func (events *requestEvents) PublishRequestApproved(request *entity.Request) error {
	msg := messaging.Message{
		Payload: request,
	}

	err := events.pub.Publish(REQUEST_APPROVED_TOPIC, msg)
	if err != nil {
		return err
	}

	return nil
}
