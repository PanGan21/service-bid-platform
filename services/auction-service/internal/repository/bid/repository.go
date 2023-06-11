package bid

import (
	"context"

	"github.com/PanGan21/pkg/entity"
)

type BidRepository interface {
	Create(ctx context.Context, bid entity.Bid) error
	FindManyByAuctionIdWithMinAmount(ctx context.Context, auctionId string) ([]entity.Bid, error)
	FindSecondMinAmountByAuctionId(ctx context.Context, auctionId string) (float64, error)
}
