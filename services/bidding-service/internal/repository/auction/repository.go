package auction

import (
	"context"

	"github.com/PanGan21/pkg/entity"
)

type AuctionRepository interface {
	Create(ctx context.Context, auction entity.Auction) error
	UpdateOne(ctx context.Context, auction entity.Auction) error
	FindOneById(ctx context.Context, id int) (entity.Auction, error)
}
