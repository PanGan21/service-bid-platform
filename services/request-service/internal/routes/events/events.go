package events

import (
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
	bidController "github.com/PanGan21/request-service/internal/routes/events/bid"
	"github.com/PanGan21/request-service/internal/service"
)

const (
	BID_CREATED_TOPIC = "bid-created"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, bidService service.BidService) {
	bidController := bidController.NewBidController(l, bidService)

	subscriber.Subscribe(BID_CREATED_TOPIC, bidController.Create)
}
