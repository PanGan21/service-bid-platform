package request

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/postgres"
)

type requestRepository struct {
	db postgres.Postgres
}

func NewRequestRepository(db postgres.Postgres) *requestRepository {
	return &requestRepository{db: db}
}

func (repo *requestRepository) Create(ctx context.Context, request entity.Request) error {
	var requestId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	const query = `
  		INSERT INTO requests (Id, CreatorId, Info, Postcode, Title, Deadline, Status, WinningBidId) 
  		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING Id;
	`

	c.QueryRow(ctx, query, request.Id, request.CreatorId, request.Info, request.Postcode, request.Title, request.Deadline, request.Status, request.WinningBidId).Scan(&requestId)
	if err != nil {
		return fmt.Errorf("RequestRepo - Create - c.Exec: %w", err)
	}

	return nil
}

func (repo *requestRepository) UpdateOne(ctx context.Context, request entity.Request) error {
	var requestId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	const query = `
  		UPDATE requests SET CreatorId=$1, Info=$2, Postcode=$3, Title=$4, Deadline=$5, Status=$6, WinningBidId=$7 WHERE Id=$8 RETURNING Id;
	`

	c.QueryRow(ctx, query, request.CreatorId, request.Info, request.Postcode, request.Title, request.Deadline, request.Status, request.WinningBidId, request.Id).Scan(&requestId)
	if err != nil {
		return fmt.Errorf("RequestRepo - Create - c.Exec: %w", err)
	}

	return nil
}
