package request

import (
	"context"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
)

type RequestRepository interface {
	Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64, status entity.RequestStatus, rejectionReason string) (int, error)
	FindOneById(ctx context.Context, id int) (entity.Request, error)
	UpdateStatusAndRejectionReasonById(ctx context.Context, id int, status entity.RequestStatus, rejectionReason string) (entity.Request, error)
	GetAllByStatus(ctx context.Context, status entity.RequestStatus, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountAllByStatus(ctx context.Context, status entity.RequestStatus) (int, error)
	GetManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string) (int, error)
}
