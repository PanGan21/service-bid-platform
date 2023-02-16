package bid

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/postgres"
)

type bidRepository struct {
	db postgres.Postgres
}

func NewBidRepository(db postgres.Postgres) *bidRepository {
	return &bidRepository{db: db}
}

func (repo *bidRepository) Create(ctx context.Context, bid entity.Bid) error {
	var bidId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	const query = `
		INSERT INTO bids (Id, Amount, CreatorId, AuctionId)
		VALUES ($1, $2, $3, $4) RETURNING Id;
	`

	err = c.QueryRow(ctx, query, bid.Id, bid.Amount, bid.CreatorId, bid.AuctionId).Scan(&bidId)
	if err != nil {
		return fmt.Errorf("AuctionRepo - Create - c.Exec: %w", err)
	}

	return nil
}

func (repo *bidRepository) FindManyByAuctionIdWithMinAmount(ctx context.Context, auctionId string) ([]entity.Bid, error) {
	var bids []entity.Bid

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return bids, err
	}
	defer c.Release()

	const query = `
		SELECT * FROM bids WHERE AuctionId=$1 AND Amount = ( SELECT MIN(Amount) FROM bids WHERE AuctionId=$1  );
	`

	rows, err := c.Query(ctx, query, auctionId)
	if err != nil {
		return bids, fmt.Errorf("AuctionRepo - FindManyByAuctionIdWithMinAmount - c.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var b entity.Bid
		err := rows.Scan(&b.Id, &b.Amount, &b.CreatorId, &b.AuctionId)
		if err != nil {
			return bids, fmt.Errorf("AuctionRepo - FindManyByAuctionIdWithMinAmount - rows.Scan: %w", err)
		}
		bids = append(bids, b)
	}

	if err := rows.Err(); err != nil {
		return bids, fmt.Errorf("AuctionRepo - FindManyByAuctionIdWithMinAmount - rows.Err: %w", err)
	}

	return bids, nil
}

func (repo *bidRepository) FindSecondMinAmountByAuctionId(ctx context.Context, auctionId string) (float64, error) {
	var secondMinAmount = 0.0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return secondMinAmount, err
	}
	defer c.Release()

	const query = `
		SELECT MIN(Amount) FROM bids WHERE AuctionId=$1 AND Amount > (SELECT MIN(Amount) From bids)
	`

	err = c.QueryRow(ctx, query, auctionId).Scan(&secondMinAmount)
	if err != nil {
		return secondMinAmount, fmt.Errorf("AuctionRepo - FindSecondMinAmountByAuctionId - c.Exec: %w", err)
	}

	return secondMinAmount, nil
}
