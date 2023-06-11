package bid

import (
	"context"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
)

type BidRepository interface {
	Create(ctx context.Context, creatorId string, auctionId int, amount float64) (int, error)
	FindOneById(ctx context.Context, id int) (entity.Bid, error)
	FindByAuctionId(ctx context.Context, auctionId int, pagination *pagination.Pagination) (*[]entity.Bid, error)
	FindByCreatorId(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Bid, error)
	CountByCreatorId(ctx context.Context, creatorId string) (int, error)
}
