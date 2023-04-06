package request

import (
	"context"
	"log"

	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

type RequestController interface {
	Create(msg messaging.Message) error
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

func (controller *requestController) Create(msg messaging.Message) error {
	request, err := entity.IsRequestType(msg.Payload)
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
		Deadline:      msg.Timestamp,
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
