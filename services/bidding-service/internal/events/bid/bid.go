package bid

import (
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/messaging"
)

type bidEvents struct {
	pub messaging.Publisher
}

const (
	BID_CREATED_TOPIC = "bid-created"
)

func NewBidEvents(pub messaging.Publisher) *bidEvents {
	return &bidEvents{pub: pub}
}

func (events *bidEvents) PublishBidCreated(request *entity.Bid) error {
	msg := messaging.Message{
		Payload: request,
	}

	err := events.pub.Publish(BID_CREATED_TOPIC, msg)
	if err != nil {
		return err
	}

	return nil
}
