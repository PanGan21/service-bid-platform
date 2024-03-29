package bid

import (
	"context"
	"log"

	"github.com/PanGan21/auction-service/internal/service"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
)

type BidController interface {
	Create(msg messaging.Message) error
}

type bidController struct {
	logger     logger.Interface
	bidService service.BidService
}

func NewBidController(logger logger.Interface, bidServ service.BidService) BidController {
	return &bidController{
		logger:     logger,
		bidService: bidServ,
	}
}

func (controller *bidController) Create(msg messaging.Message) error {
	bid, err := entity.IsBidType(msg.Payload)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	err = controller.bidService.Create(context.Background(), bid)
	if err != nil {
		controller.logger.Error(err)
		log.Fatal(err)
	}

	return nil
}
