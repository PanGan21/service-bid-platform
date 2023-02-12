package auction

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PanGan21/auction-service/internal/service"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/pagination"
	"github.com/PanGan21/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuctionController interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	CountAll(c *gin.Context)
	GetOwn(c *gin.Context)
	CountOwn(c *gin.Context)
	UpdateWinnerByAuctionId(c *gin.Context)
	GetOpenPastDeadline(c *gin.Context)
	CountOpenPastDeadline(c *gin.Context)
	UpdateStatus(c *gin.Context)
	GetByStatus(c *gin.Context)
	CountByStatus(c *gin.Context)
	GetOwnAssignedByStatuses(c *gin.Context)
	CountOwnAssignedByStatuses(c *gin.Context)
}

type auctionController struct {
	logger         logger.Interface
	auctionService service.AuctionService
	bidService     service.BidService
}

func NewAuctionController(logger logger.Interface, auctionServ service.AuctionService, bidServ service.BidService) AuctionController {
	return &auctionController{
		logger:         logger,
		auctionService: auctionServ,
		bidService:     bidServ,
	}
}

type AuctionData struct {
	Title    string `json:"Title"`
	Postcode string `json:"Postcode"`
	Info     string `json:"Info"`
	Deadline int64  `json:"Deadline"`
}

func (controller *auctionController) Create(c *gin.Context) {
	var auctionData AuctionData
	if err := c.BindJSON(&auctionData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	auction, err := controller.auctionService.Create(context.Background(), userId.(string), auctionData.Info, auctionData.Postcode, auctionData.Title, auctionData.Deadline)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Creation failed"})
		return
	}

	c.JSON(http.StatusOK, auction)
}

func (controller *auctionController) GetAll(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)

	auctions, err := controller.auctionService.GetAll(context.Background(), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, auctions)
}

func (controller *auctionController) CountAll(c *gin.Context) {
	count, err := controller.auctionService.CountAll(context.Background())
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}

func (controller *auctionController) GetOwn(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	pagination := pagination.GeneratePaginationFromRequest(c)

	auctions, err := controller.auctionService.GetOwn(context.Background(), userId.(string), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, auctions)
}

func (controller *auctionController) CountOwn(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	count, err := controller.auctionService.CountOwn(context.Background(), userId.(string))
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}

func (controller *auctionController) UpdateWinnerByAuctionId(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("auctionId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	auction, err := controller.auctionService.GetById(context.Background(), id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found"})
		return
	}

	isAllowedToResolve := controller.auctionService.IsAllowedToResolve(context.Background(), auction)
	if !isAllowedToResolve {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auction not allowed to be resolved"})
		return
	}

	winnignBid, err := controller.bidService.FindWinningBidByAuctionId(context.Background(), idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find winning bid"})
		return
	}

	_, err = controller.auctionService.UpdateWinningBid(context.Background(), auction, strconv.Itoa(winnignBid.Id))
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not update auction"})
		return
	}

	c.JSON(http.StatusOK, winnignBid)
}

func (controller *auctionController) GetOpenPastDeadline(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)

	auctions, err := controller.auctionService.GetAllOpenPastDeadline(context.Background(), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (controller *auctionController) CountOpenPastDeadline(c *gin.Context) {
	count, err := controller.auctionService.CountAllOpenPastDeadline(context.Background())
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}

type UpdateStatusData struct {
	Status entity.AuctionStatus `json:"Status"`
}

func (controller *auctionController) UpdateStatus(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("auctionId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	var updateStatusData UpdateStatusData
	if err := c.BindJSON(&updateStatusData); err != nil || (updateStatusData.Status != entity.InProgress && updateStatusData.Status != entity.Closed && updateStatusData.Status != entity.Open) {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	updatedAuction, err := controller.auctionService.UpdateStatusByAuctionId(context.Background(), updateStatusData.Status, id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, updatedAuction)
}

func (controller *auctionController) GetByStatus(c *gin.Context) {
	statusParam := c.Request.URL.Query().Get("status")
	allowedStatuses := []string{string(entity.Open), string(entity.New), string(entity.InProgress), string(entity.Assigned), string(entity.Closed)}

	if !utils.Contains(allowedStatuses, statusParam) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	status := entity.AuctionStatus(statusParam)

	pagination := pagination.GeneratePaginationFromRequest(c)

	auctions, err := controller.auctionService.GetAllByStatus(context.Background(), status, &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (controller *auctionController) CountByStatus(c *gin.Context) {
	statusParam := c.Request.URL.Query().Get("status")
	allowedStatuses := []string{string(entity.Open), string(entity.New), string(entity.InProgress), string(entity.Assigned), string(entity.Closed)}

	if !utils.Contains(allowedStatuses, statusParam) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	status := entity.AuctionStatus(statusParam)

	count, err := controller.auctionService.CountAllByStatus(context.Background(), status)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)
}

func (controller *auctionController) GetOwnAssignedByStatuses(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	pagination := pagination.GeneratePaginationFromRequest(c)

	statuses := []entity.AuctionStatus{entity.Assigned, entity.InProgress}

	populatedAuctions, err := controller.auctionService.GetOwnAssignedByStatuses(context.Background(), statuses, userId.(string), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, populatedAuctions)
}

func (controller *auctionController) CountOwnAssignedByStatuses(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	statuses := []entity.AuctionStatus{entity.Assigned, entity.InProgress}

	count, err := controller.auctionService.CountOwnAssignedByStatuses(context.Background(), statuses, userId.(string))
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)
}
