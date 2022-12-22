package bid

import (
	"context"
	"net/http"

	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/logger"
	"github.com/gin-gonic/gin"
)

type BidController interface {
	Create(c *gin.Context)
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

type BidData struct {
	Amount    float64 `json:"Amount"`
	RequestId int     `json:"RequestId"`
}

func (controller *bidController) Create(c *gin.Context) {
	var bidData BidData
	if err := c.BindJSON(&bidData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	bid, err := controller.bidService.Create(context.Background(), userId.(string), bidData.RequestId, bidData.Amount)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Creation failed"})
		return
	}

	c.JSON(http.StatusOK, bid)
}
