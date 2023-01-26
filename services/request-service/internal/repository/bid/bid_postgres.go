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

	c.QueryRow(ctx, query, bid.Id, bid.Amount, bid.CreatorId, bid.RequestId).Scan(&bidId)
	if err != nil {
		return fmt.Errorf("RequestRepo - Create - c.Exec: %w", err)
	}

	return nil
}
