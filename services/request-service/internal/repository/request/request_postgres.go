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

func (repo *requestRepository) Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64, status entity.RequestStatus, rejectionReason string) (int, error) {
	var requestId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return requestId, err
	}
	defer c.Release()

	const query = `
  		INSERT INTO requests (CreatorId, Info, Postcode, Title, Deadline, Status, RejectionReason) 
  		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING Id;
	`

	c.QueryRow(ctx, query, creatorId, info, postcode, title, deadline, status, rejectionReason).Scan(&requestId)
	if err != nil {
		return requestId, fmt.Errorf("RequestRepo - Create - c.QueryRow: %w", err)
	}

	return requestId, nil
}

func (repo *requestRepository) FindOneById(ctx context.Context, id int) (entity.Request, error) {
	var request entity.Request

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return request, err
	}
	defer c.Release()

	const query = `
		SELECT * FROM requests WHERE Id=$1;
	`

	err = c.QueryRow(ctx, query, id).Scan(&request.Id, &request.Title, &request.Postcode, &request.Info, &request.CreatorId, &request.Deadline, &request.Status, &request.RejectionReason)
	if err != nil {
		return request, fmt.Errorf("RequestRepo - FindOneById - c.QueryRow: %w", err)
	}

	return request, nil
}
