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
  		INSERT INTO requests (Id, CreatorId, Info, Postcode, Title, Deadline) 
  		VALUES ($1, $2, $3, $4, $5, $6) RETURNING Id;
	`

	c.QueryRow(ctx, query, request.Id, request.CreatorId, request.Info, request.Postcode, request.Title, request.Deadline).Scan(&requestId)
	if err != nil {
		return fmt.Errorf("RequestRepo - Create - c.Exec: %w", err)
	}

	return nil
}
