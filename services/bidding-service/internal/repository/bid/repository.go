package bid

import (
	"context"

	"github.com/PanGan21/pkg/entity"
)

type BidRepository interface {
	Create(ctx context.Context, creatorId string, requestId int, amount float64) (int, error)
	FindOneById(ctx context.Context, id int) (entity.Bid, error)
}
