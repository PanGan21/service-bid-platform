package auction

import (
	"context"
	"log"

	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
)

type AuctionController interface {
	Create(payload interface{}) error
	Update(payload interface{}) error
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

func (controller *auctionController) Create(payload interface{}) error {
	auction, err := entity.IsAuctionType(payload)
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

func (controller *auctionController) Update(payload interface{}) error {
	auction, err := entity.IsAuctionType(payload)
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
