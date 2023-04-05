package request

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/pagination"
	"github.com/PanGan21/pkg/utils"
	"github.com/PanGan21/request-service/internal/service"
	"github.com/gin-gonic/gin"
)

type RequestController interface {
	Create(c *gin.Context)
	RejectRequest(c *gin.Context)
	GetByStatus(c *gin.Context)
	CountByStatus(c *gin.Context)
	GetOwnByStatus(c *gin.Context)
	CountOwnByStatus(c *gin.Context)
	Approve(c *gin.Context)
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

type RequestData struct {
	Title    string `json:"Title"`
	Postcode string `json:"Postcode"`
	Info     string `json:"Info"`
}

func (controller *requestController) Create(c *gin.Context) {
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	request, err := controller.requestService.Create(context.Background(), user.(entity.PublicUser).Id, requestData.Info, requestData.Postcode, requestData.Title)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Creation failed"})
		return
	}

	c.JSON(http.StatusOK, request)
}

type RejectRequestData struct {
	RejectionReason string
}

func (controller *requestController) RejectRequest(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("requestId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	var rejectRequestData RejectRequestData
	if err := c.BindJSON(&rejectRequestData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	request, err := controller.requestService.RejectRequest(context.Background(), rejectRequestData.RejectionReason, id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}

func (controller *requestController) GetByStatus(c *gin.Context) {
	statusParam := c.Request.URL.Query().Get("status")
	allowedStatuses := []string{string(entity.NewRequest), string(entity.ApprovedRequest), string(entity.RejectedRequest)}

	if !utils.Contains(allowedStatuses, statusParam) {
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
	allowedStatuses := []string{string(entity.Open), string(entity.NewRequest), string(entity.RejectedRequest)}

	if !utils.Contains(allowedStatuses, statusParam) {
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

func (controller *requestController) GetOwnByStatus(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	statusParam := c.Request.URL.Query().Get("status")
	allowedStatuses := []string{string(entity.Open), string(entity.NewRequest), string(entity.RejectedRequest)}

	if !utils.Contains(allowedStatuses, statusParam) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	status := entity.RequestStatus(statusParam)

	pagination := pagination.GeneratePaginationFromRequest(c)

	requests, err := controller.requestService.GetManyByStatusByUserId(context.Background(), status, user.(entity.PublicUser).Id, &pagination)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, requests)
}

func (controller *requestController) CountOwnByStatus(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	statusParam := c.Request.URL.Query().Get("status")
	allowedStatuses := []string{string(entity.Open), string(entity.NewRequest), string(entity.RejectedRequest)}

	if !utils.Contains(allowedStatuses, statusParam) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	status := entity.RequestStatus(statusParam)

	count, err := controller.requestService.CountManyByStatusByUserId(context.Background(), status, user.(entity.PublicUser).Id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, count)
}

func (controller *requestController) Approve(c *gin.Context) {
	idParam := c.Request.URL.Query().Get("requestId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	request, err := controller.requestService.ApproveRequestById(context.Background(), id)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}
