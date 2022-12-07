package pagination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	// Default
	limit := 2
	page := 1
	asc := true

	var err error

	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, err = strconv.Atoi(queryValue)
		case "page":
			page, err = strconv.Atoi(queryValue)
		case "asc":
			asc, err = strconv.ParseBool(queryValue)

		}
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
	}

	return Pagination{
		Limit: limit,
		Page:  page,
		Asc:   asc,
	}
}
