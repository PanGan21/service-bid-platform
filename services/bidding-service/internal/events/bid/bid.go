package bid

import (
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/messaging"
)

type bidEvents struct {
	pub messaging.Publisher
}

func NewBidEvents(pub messaging.Publisher) *bidEvents {
	return &bidEvents{pub: pub}
}

func (events *bidEvents) PublishBidCreated(topic string, request *entity.Bid) error {
	msg := messaging.Message{
		Payload: request,
	}

	err := events.pub.Publish(topic, msg)
	if err != nil {
		return err
	}

	return nil
}
