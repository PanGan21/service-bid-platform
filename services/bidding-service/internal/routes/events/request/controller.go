package request

import (
	"context"
	"log"

	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
)

type RequestController interface {
	Create(payload interface{}) error
	Update(payload interface{}) error
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
	request, err := entity.IsRequestType(payload)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	err = controller.requestService.Create(context.Background(), request)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	return nil
}

func (controller *requestController) Update(payload interface{}) error {
	request, err := entity.IsRequestType(payload)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	err = controller.requestService.UpdateOne(context.Background(), request)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	return nil
}
