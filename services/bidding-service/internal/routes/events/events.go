package events

import (
	requestController "github.com/PanGan21/bidding-service/internal/routes/events/request"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

const (
	REQUEST_CREATED_TOPIC = "request-created"
	REQUEST_UPDATED_TOPIC = "request-updated"
)

func NewEventsClient(subscriber messaging.Subscriber, l logger.Interface, requestService service.RequestService) {
	requestController := requestController.NewRequestController(l, requestService)

	go subscriber.Subscribe(REQUEST_CREATED_TOPIC, requestController.Create)
	go subscriber.Subscribe(REQUEST_UPDATED_TOPIC, requestController.Update)
}
