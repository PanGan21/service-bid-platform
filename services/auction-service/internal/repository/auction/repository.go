package auction

import (
	"context"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
)

type AuctionRepository interface {
	Create(ctx context.Context, auction entity.Auction) (int, error)
	GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountAll(ctx context.Context) (int, error)
	FindOneById(ctx context.Context, id int) (entity.Auction, error)
	FindByCreatorId(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountByCreatorId(ctx context.Context, creatorId string) (int, error)
	UpdateWinningBidIdAndStatusById(ctx context.Context, id int, winningBidId string, status entity.AuctionStatus, winnerId string, winningAmount float64) (entity.Auction, error)
	GetAllOpenPastTime(ctx context.Context, timestamp int64, pagination *pagination.Pagination) (*[]entity.ExtendedAuction, error)
	CountAllOpenPastTime(ctx context.Context, timestamp int64) (int, error)
	UpdateStatusByAuctionId(ctx context.Context, status entity.AuctionStatus, id int) (entity.Auction, error)
	GetAllByStatus(ctx context.Context, status entity.AuctionStatus, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountAllByStatus(ctx context.Context, status entity.AuctionStatus) (int, error)
	GetOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string) (int, error)
	FindByCreatorIdAndStatus(ctx context.Context, creatorId string, status entity.AuctionStatus, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountByCreatorIdAndStatus(ctx context.Context, creatorId string, status entity.AuctionStatus) (int, error)
	UpdateDeadlineByAuctionId(ctx context.Context, deadline int64, id int) (entity.Auction, error)
}
