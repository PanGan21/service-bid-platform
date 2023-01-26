package bid

import (
	"context"

	"github.com/PanGan21/pkg/entity"
)

type BidRepository interface {
	Create(ctx context.Context, bid entity.Bid) error
}
