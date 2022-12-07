package request

import (
	"context"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
)

type RequestRepository interface {
	Create(ctx context.Context, request *entity.Request) error
	GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Request, error)
	FindByCreatorId(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Request, error)
}
