package bid

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type BidController interface {
	Create(c *gin.Context)
	GetOneById(c *gin.Context)
	GetManyByAuctionId(c *gin.Context)
	GetOwn(c *gin.Context)
	CountOwn(c *gin.Context)
}

type bidController struct {
	logger        logger.Interface
	bidService    service.BidService
	auctionServce service.AuctionService
}

func NewBidController(logger logger.Interface, bidServ service.BidService, reqServ service.AuctionService) BidController {
	return &bidController{
		logger:        logger,
		bidService:    bidServ,
		auctionServce: reqServ,
	}
}

type BidData struct {
	Amount    float64 `json:"Amount"`
	AuctionId int     `json:"AuctionId"`
}

func (controller *bidController) Create(c *gin.Context) {
	var bidData BidData
	if err := c.BindJSON(&bidData); err != nil {
		controller.logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
		return
	}

	isAllowedToBeCreated := controller.auctionServce.IsOpenToBidByAuctionId(context.Background(), bidData.AuctionId)
	if !isAllowedToBeCreated {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auction doesn't receive bids"})
		return
	}

	bid, err := controller.bidService.Create(context.Background(), user.(entity.PublicUser).Id, bidData.AuctionId, bidData.Amount)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Creation failed"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func (controller *bidController) GetOneById(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	bid, err := controller.bidService.FindOneById(c.Request.Context(), id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find bid"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func (controller *bidController) GetManyByAuctionId(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("auctionId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	pagination := pagination.GeneratePaginationFromRequest(c)

	bids, err := controller.bidService.GetManyByAuctionId(c.Request.Context(), id, &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, bids)
}

func (controller *bidController) GetOwn(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	pagination := pagination.GeneratePaginationFromRequest(c)

	bids, err := controller.bidService.GetOwn(c.Request.Context(), user.(entity.PublicUser).Id, &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, bids)
}

func (controller *bidController) CountOwn(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	count, err := controller.bidService.CountOwn(context.Background(), user.(entity.PublicUser).Id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}
