package auction

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/postgres"
)

type auctionRepository struct {
	db postgres.Postgres
}

func NewAuctionRepository(db postgres.Postgres) *auctionRepository {
	return &auctionRepository{db: db}
}

func (repo *auctionRepository) Create(ctx context.Context, auction entity.Auction) error {
	var auctionId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	const query = `
  		INSERT INTO auctions (Id, CreatorId, Info, Postcode, Title, Deadline, Status, WinningBidId, RejectionReason) 
  		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING Id;
	`

	c.QueryRow(ctx, query, auction.Id, auction.CreatorId, auction.Info, auction.Postcode, auction.Title, auction.Deadline, auction.Status, auction.WinningBidId, auction.RejectionReason).Scan(&auctionId)
	if err != nil {
		return fmt.Errorf("AuctionRepo - Create - c.Exec: %w", err)
	}

	return nil
}

func (repo *auctionRepository) UpdateOne(ctx context.Context, auction entity.Auction) error {
	var auctionId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	const query = `
  		UPDATE auctions SET CreatorId=$1, Info=$2, Postcode=$3, Title=$4, Deadline=$5, Status=$6, WinningBidId=$7, RejectionReason=$8 WHERE Id=$9 RETURNING Id;
	`

	c.QueryRow(ctx, query, auction.CreatorId, auction.Info, auction.Postcode, auction.Title, auction.Deadline, auction.Status, auction.WinningBidId, auction.RejectionReason, auction.Id).Scan(&auctionId)
	if err != nil {
		return fmt.Errorf("AuctionRepo - Create - c.Exec: %w", err)
	}

	return nil
}

func (repo *auctionRepository) FindOneById(ctx context.Context, id int) (entity.Auction, error) {
	var auction entity.Auction

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return auction, err
	}
	defer c.Release()

	const query = `
		SELECT * FROM auctions WHERE Id=$1;
	`

	err = c.QueryRow(ctx, query, id).Scan(&auction.Id, &auction.Title, &auction.Postcode, &auction.Info, &auction.CreatorId, &auction.Deadline, &auction.Status, &auction.WinningBidId, &auction.RejectionReason)
	if err != nil {
		return auction, fmt.Errorf("AuctionRepo - FindOneById - c.Exec: %w", err)
	}

	return auction, nil
}
