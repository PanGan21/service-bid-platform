package events

import (
	bidController "github.com/PanGan21/auction-service/internal/routes/events/bid"
	"github.com/PanGan21/auction-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

const (
	BID_CREATED_TOPIC = "bid-created"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, bidService service.BidService) {
	bidController := bidController.NewBidController(l, bidService)

	go subscriber.Subscribe(BID_CREATED_TOPIC, bidController.Create)
}
