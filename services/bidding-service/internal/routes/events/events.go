package events

import (
	auctionController "github.com/PanGan21/bidding-service/internal/routes/events/auction"
	requestController "github.com/PanGan21/bidding-service/internal/routes/events/request"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

const (
	REQUEST_APPROVED_TOPIC = "request-approved"
	AUCTION_UPDATED_TOPIC  = "auction-updated"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, auctionService service.AuctionService) {
	auctionController := auctionController.NewAuctionController(l, auctionService)
	requestController := requestController.NewRequestController(l, auctionService)

	go subscriber.Subscribe(REQUEST_APPROVED_TOPIC, requestController.Create)
	go subscriber.Subscribe(AUCTION_UPDATED_TOPIC, auctionController.Update)
}
