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
}

type requestController struct {
	logger         logger.Interface
	auctionService service.AuctionService
}

func NewRequestController(logger logger.Interface, auctionServ service.AuctionService) RequestController {
	return &requestController{
		logger:         logger,
		auctionService: auctionServ,
	}
}

func (controller *requestController) Create(payload interface{}) error {
	request, err := entity.IsRequestType(payload)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	newAuction := entity.Auction{
		Id:            request.Id,
		Title:         request.Title,
		Postcode:      request.Postcode,
		Info:          request.Info,
		CreatorId:     request.CreatorId,
		Deadline:      request.Deadline,
		Status:        entity.Open,
		WinningBidId:  "",
		WinnerId:      "",
		WinningAmount: 0.0,
	}

	err = controller.auctionService.Create(context.Background(), newAuction)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	return nil
}
