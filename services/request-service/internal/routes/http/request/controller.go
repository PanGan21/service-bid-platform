package request

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/pagination"
	"github.com/PanGan21/request-service/internal/service"
	"github.com/gin-gonic/gin"
)

type RequestController interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	CountAll(c *gin.Context)
	GetOwn(c *gin.Context)
	CountOwn(c *gin.Context)
	UpdateWinnerByRequestId(c *gin.Context)
	GetOpenPastDeadline(c *gin.Context)
	CountOpenPastDeadline(c *gin.Context)
	UpdateStatus(c *gin.Context)
	GetByStatus(c *gin.Context)
	CountByStatus(c *gin.Context)
	GetOwnAssignedByStatuses(c *gin.Context)
}

type requestController struct {
	logger         logger.Interface
	requestService service.RequestService
	bidService     service.BidService
}

func NewRequestController(logger logger.Interface, requestServ service.RequestService, bidServ service.BidService) RequestController {
	return &requestController{
		logger:         logger,
		requestService: requestServ,
		bidService:     bidServ,
	}
}

type RequestData struct {
	Title    string `json:"Title"`
	Postcode string `json:"Postcode"`
	Info     string `json:"Info"`
	Deadline int64  `json:"Deadline"`
}

func (controller *requestController) Create(c *gin.Context) {
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	request, err := controller.requestService.Create(context.Background(), userId.(string), requestData.Info, requestData.Postcode, requestData.Title, requestData.Deadline)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Creation failed"})
		return
	}

	c.JSON(http.StatusOK, request)
}

func (controller *requestController) GetAll(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)

	requests, err := controller.requestService.GetAll(context.Background(), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, requests)
}

func (controller *requestController) CountAll(c *gin.Context) {
	count, err := controller.requestService.CountAll(context.Background())
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}

func (controller *requestController) GetOwn(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	pagination := pagination.GeneratePaginationFromRequest(c)

	requests, err := controller.requestService.GetOwn(context.Background(), userId.(string), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, requests)
}

func (controller *requestController) CountOwn(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	count, err := controller.requestService.CountOwn(context.Background(), userId.(string))
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}

func (controller *requestController) UpdateWinnerByRequestId(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("requestId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	request, err := controller.requestService.GetById(context.Background(), id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	isAllowedToResolve := controller.requestService.IsAllowedToResolve(context.Background(), request)
	if !isAllowedToResolve {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Request not allowed to be resolved"})
		return
	}

	winnignBid, err := controller.bidService.FindWinningBidByRequestId(context.Background(), idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find winning bid"})
		return
	}

	_, err = controller.requestService.UpdateWinningBid(context.Background(), request, strconv.Itoa(winnignBid.Id))
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not update request"})
		return
	}

	c.JSON(http.StatusOK, winnignBid)
}

func (controller *requestController) GetOpenPastDeadline(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)

	requests, err := controller.requestService.GetAllOpenPastDeadline(context.Background(), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (controller *requestController) CountOpenPastDeadline(c *gin.Context) {
	count, err := controller.requestService.CountAllOpenPastDeadline(context.Background())
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}

type UpdateStatusData struct {
	Status entity.RequestStatus `json:"Status"`
}

func (controller *requestController) UpdateStatus(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("requestId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	var updateStatusData UpdateStatusData
	if err := c.BindJSON(&updateStatusData); err != nil || (updateStatusData.Status != entity.InProgress && updateStatusData.Status != entity.Closed) {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	updatedRequest, err := controller.requestService.UpdateStatusByRequestId(context.Background(), updateStatusData.Status, id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, updatedRequest)
}

func (controller *requestController) GetByStatus(c *gin.Context) {
	statusParam := c.Request.URL.Query().Get("status")
	if statusParam == "" || (statusParam != string(entity.Open) && statusParam != string(entity.Assigned) && statusParam != string(entity.Closed) && statusParam != string(entity.InProgress)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	status := entity.RequestStatus(statusParam)

	pagination := pagination.GeneratePaginationFromRequest(c)

	requests, err := controller.requestService.GetAllByStatus(context.Background(), status, &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (controller *requestController) CountByStatus(c *gin.Context) {
	statusParam := c.Request.URL.Query().Get("status")
	if statusParam == "" || (statusParam != string(entity.Open) && statusParam != string(entity.Assigned) && statusParam != string(entity.Closed) && statusParam != string(entity.InProgress)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	status := entity.RequestStatus(statusParam)

	count, err := controller.requestService.CountAllByStatus(context.Background(), status)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)
}

func (controller *requestController) GetOwnAssignedByStatuses(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	pagination := pagination.GeneratePaginationFromRequest(c)

	statuses := []entity.RequestStatus{entity.Assigned, entity.InProgress}

	populatedRequests, err := controller.requestService.GetOwnAssignedByStatuses(context.Background(), statuses, userId.(string), &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, populatedRequests)
}
