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
		INSERT INTO bids (Id, Amount, CreatorId, RequestId)
		VALUES ($1, $2, $3, $4) RETURNING Id;
	`

	err = c.QueryRow(ctx, query, bid.Id, bid.Amount, bid.CreatorId, bid.RequestId).Scan(&bidId)
	if err != nil {
		return fmt.Errorf("RequestRepo - Create - c.Exec: %w", err)
	}

	return nil
}

func (repo *bidRepository) FindManyByRequestIdWithMinAmount(ctx context.Context, requestId string) ([]entity.Bid, error) {
	var bids []entity.Bid

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return bids, err
	}
	defer c.Release()

	const query = `
		SELECT * FROM bids WHERE RequestId = $1 AND Amount = ( SELECT MIN(Amount) FROM bids );
	`

	rows, err := c.Query(ctx, query, requestId)
	if err != nil {
		return bids, fmt.Errorf("RequestRepo - FindManyByRequestIdWithMinAmount - c.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var b entity.Bid
		err := rows.Scan(&b.Id, &b.Amount, &b.CreatorId, &b.RequestId)
		if err != nil {
			return bids, fmt.Errorf("RequestRepo - FindManyByRequestIdWithMinAmount - rows.Scan: %w", err)
		}
		bids = append(bids, b)
	}

	if err := rows.Err(); err != nil {
		return bids, fmt.Errorf("RequestRepo - FindManyByRequestIdWithMinAmount - rows.Err: %w", err)
	}

	return bids, nil
}
