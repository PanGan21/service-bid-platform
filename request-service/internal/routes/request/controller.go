package request

import (
	"fmt"
	"net/http"

	"github.com/PanGan21/packages/logger"
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
	Title    string `json:"title"`
	Postcode string `json:"postcode"`
	Info     string `json:"info"`
	Deadline int64  `json:"deadline"`
}

func (controller *requestController) Create(c *gin.Context) {
	// TODO: validate
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	fmt.Println("requestData", requestData)

	_, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator does not exist; Authentication error"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created request"})
}
