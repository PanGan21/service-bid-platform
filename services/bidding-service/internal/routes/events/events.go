package events

import (
	auctionController "github.com/PanGan21/bidding-service/internal/routes/events/auction"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

const (
	AUCTION_CREATED_TOPIC = "auction-created"
	AUCTION_UPDATED_TOPIC = "auction-updated"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, auctionService service.AuctionService) {
	auctionController := auctionController.NewAuctionController(l, auctionService)

	go subscriber.Subscribe(AUCTION_CREATED_TOPIC, auctionController.Create)
	go subscriber.Subscribe(AUCTION_UPDATED_TOPIC, auctionController.Update)
}
