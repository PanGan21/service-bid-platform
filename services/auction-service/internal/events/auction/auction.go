package auction

import (
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/messaging"
)

type auctionEvents struct {
	pub messaging.Publisher
}

const (
	AUCTION_UPDATED_TOPIC = "auction-updated"
)

func NewAuctionEvents(pub messaging.Publisher) *auctionEvents {
	return &auctionEvents{pub: pub}
}

func (events *auctionEvents) PublishAuctionUpdated(auction *entity.Auction) error {
	msg := messaging.Message{
		Payload: auction,
	}

	err := events.pub.Publish(AUCTION_UPDATED_TOPIC, msg)
	if err != nil {
		return err
	}

	return nil
}
