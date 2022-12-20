package request

import (
	"context"
	"log"

	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/logger"
)

type RequestController interface {
	Create(payload interface{}) error
}

type requestController struct {
	logger         logger.Interface
	requestService service.RequestService
}

func NewRequestController(logger logger.Interface, requestServ service.RequestService) RequestController {
	return &requestController{
		logger:         logger,
		requestService: requestServ,
	}
}

func (controller *requestController) Create(payload interface{}) error {
	// request, ok := payload.(entity.Request)
	creatorId, err := payload["creatorId"].(string)
	// if !ok {
	// 	fmt.Println("request", request)
	// 	controller.logger.Error("incorrect payload ", payload)
	// 	log.Fatal("incorrect payload ", payload)
	// }

	err := controller.requestService.Create(context.Background(), request)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	return nil
}
