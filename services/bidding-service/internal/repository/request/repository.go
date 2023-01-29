package request

import (
	"context"

	"github.com/PanGan21/pkg/entity"
)

type RequestRepository interface {
	Create(ctx context.Context, request entity.Request) error
	UpdateOne(ctx context.Context, request entity.Request) error
}
