package events

import (
	requestController "github.com/PanGan21/bidding-service/internal/routes/events/request"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

const (
	REQUEST_CREATED_TOPIC = "request-created"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, requestService service.RequestService) {
	requestController := requestController.NewRequestController(l, requestService)

	subscriber.Subscribe(REQUEST_CREATED_TOPIC, requestController.Create)
}
