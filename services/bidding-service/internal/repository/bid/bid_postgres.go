package bid

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
	"github.com/PanGan21/pkg/postgres"
)

type bidRepository struct {
	db postgres.Postgres
}

func NewBidRepository(db postgres.Postgres) *bidRepository {
	return &bidRepository{db: db}
}

func (repo *bidRepository) Create(ctx context.Context, creatorId string, requestId int, amount float64) (int, error) {
	var bidId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return requestId, err
	}
	defer c.Release()

	const query = `
		INSERT INTO bids (Amount, CreatorId, RequestId)
		VALUES ($1, $2, $3) RETURNING Id;
	`

	c.QueryRow(ctx, query, amount, creatorId, requestId).Scan(&bidId)
	if err != nil {
		return bidId, fmt.Errorf("BidRepo - Create - c.QueryRow: %w", err)
	}

	return bidId, nil
}

func (repo *bidRepository) FindOneById(ctx context.Context, id int) (entity.Bid, error) {
	var bid entity.Bid

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return bid, err
	}
	defer c.Release()

	const query = `
		SELECT * FROM bids WHERE Id=$1;
	`

	err = c.QueryRow(ctx, query, id).Scan(&bid.Id, &bid.Amount, &bid.CreatorId, &bid.RequestId)
	if err != nil {
		return bid, fmt.Errorf("BidRepo - FindOneById - c.QueryRow: %w", err)
	}

	return bid, nil
}

func (repo *bidRepository) FindByRequestId(ctx context.Context, requestId int, pagination *pagination.Pagination) (*[]entity.Bid, error) {
	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()

	offset := (pagination.Page - 1) * pagination.Limit

	order := "asc"
	if !pagination.Asc {
		order = "desc"
	}

	query := fmt.Sprintf("SELECT * FROM bids WHERE RequestId=$1 ORDER BY Id %s LIMIT $2 OFFSET $3;", order)

	rows, err := c.Query(ctx, query, requestId, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("BidRepo - FindByRequestId - c.Query: %w", err)
	}
	defer rows.Close()

	var bids []entity.Bid

	for rows.Next() {
		var b entity.Bid
		err := rows.Scan(&b.Id, &b.Amount, &b.CreatorId, &b.RequestId)
		if err != nil {
			return nil, fmt.Errorf("BidRepo - FindByRequestId - rows.Scan: %w", err)
		}
		bids = append(bids, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("BidRepo - FindByRequestId - rows.Err: %w", err)
	}

	return &bids, nil
}
