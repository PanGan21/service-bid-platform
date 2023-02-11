package events

import (
	auctionController "github.com/PanGan21/bidding-service/internal/routes/events/auction"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

const (
	REQUEST_CREATED_TOPIC = "auction-created"
	REQUEST_UPDATED_TOPIC = "auction-updated"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, auctionService service.AuctionService) {
	auctionController := auctionController.NewAuctionController(l, auctionService)

	go subscriber.Subscribe(REQUEST_CREATED_TOPIC, auctionController.Create)
	go subscriber.Subscribe(REQUEST_UPDATED_TOPIC, auctionController.Update)
}
