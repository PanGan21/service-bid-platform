package request

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/request-service/internal/service"
	"github.com/gin-gonic/gin"
)

type RequestController interface {
	Create(c *gin.Context)
	RejectRequest(c *gin.Context)
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
	Deadline int64  `json:"Deadline"`
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

	request, err := controller.requestService.Create(context.Background(), user.(entity.PublicUser).Id, requestData.Info, requestData.Postcode, requestData.Title, requestData.Deadline)
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
