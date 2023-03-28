package request

import (
	"context"
	"net/http"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/request-service/internal/service"
	"github.com/gin-gonic/gin"
)

type RequestController interface {
	Create(c *gin.Context)
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
