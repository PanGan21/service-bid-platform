package request

import (
	"context"

	"github.com/PanGan21/request-service/internal/entity"
)

type RequestRepository interface {
	Create(ctx context.Context, request *entity.Request) error
}
