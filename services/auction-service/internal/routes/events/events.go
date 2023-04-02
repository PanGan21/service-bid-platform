package events

import (
	bidController "github.com/PanGan21/auction-service/internal/routes/events/bid"
	requestController "github.com/PanGan21/auction-service/internal/routes/events/request"
	"github.com/PanGan21/auction-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

const (
	BID_CREATED_TOPIC      = "bid-created"
	REQUEST_APPROVED_TOPIC = "request-approved"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, bidService service.BidService, auctionService service.AuctionService) {
	bidController := bidController.NewBidController(l, bidService)
	requestController := requestController.NewRequestController(l, auctionService)

	go subscriber.Subscribe(BID_CREATED_TOPIC, bidController.Create)
	go subscriber.Subscribe(REQUEST_APPROVED_TOPIC, requestController.Create)
}
