package auction

import (
	"context"
	"log"

	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

type AuctionController interface {
	Create(msg messaging.Message) error
	Update(msg messaging.Message) error
}

type auctionController struct {
	logger         logger.Interface
	auctionService service.AuctionService
}

func NewAuctionController(logger logger.Interface, auctionServ service.AuctionService) AuctionController {
	return &auctionController{
		logger:         logger,
		auctionService: auctionServ,
	}
}

func (controller *auctionController) Create(msg messaging.Message) error {
	auction, err := entity.IsAuctionType(msg.Payload)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	err = controller.auctionService.Create(context.Background(), auction)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	return nil
}

func (controller *auctionController) Update(msg messaging.Message) error {
	auction, err := entity.IsAuctionType(msg.Payload)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	err = controller.auctionService.UpdateOne(context.Background(), auction)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	return nil
}
