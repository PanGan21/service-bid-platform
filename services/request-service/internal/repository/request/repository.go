package request

import (
	"context"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
)

type RequestRepository interface {
	Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64) (int, error)
	GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Request, error)
	FindOneById(ctx context.Context, id int) (entity.Request, error)
	FindByCreatorId(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountByCreatorId(ctx context.Context, creatorId string) (int, error)
}
