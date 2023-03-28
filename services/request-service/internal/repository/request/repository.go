package request

import (
	"context"

	"github.com/PanGan21/pkg/entity"
)

type RequestRepository interface {
	Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64, status entity.RequestStatus, rejectionReason string) (int, error)
	FindOneById(ctx context.Context, id int) (entity.Request, error)
	UpdateStatusAndRejectionReasonById(ctx context.Context, id int, status entity.RequestStatus, rejectionReason string) (entity.Request, error)
}
