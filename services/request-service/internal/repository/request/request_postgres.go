package request

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
	"github.com/PanGan21/pkg/postgres"
)

type requestRepository struct {
	db postgres.Postgres
}

func NewRequestRepository(db postgres.Postgres) *requestRepository {
	return &requestRepository{db: db}
}

func (repo *requestRepository) Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64) (int, error) {
	var requestId int
	var defaultStatus = entity.Open
	var defaultWinnigBidId = ""

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return requestId, err
	}
	defer c.Release()

	const query = `
  		INSERT INTO requests (CreatorId, Info, Postcode, Title, Deadline, Status, WinningBidId) 
  		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING Id;
	`

	c.QueryRow(ctx, query, creatorId, info, postcode, title, deadline, defaultStatus, defaultWinnigBidId).Scan(&requestId)
	if err != nil {
		return requestId, fmt.Errorf("RequestRepo - Create - c.QueryRow: %w", err)
	}

	return requestId, nil
}

func (repo *requestRepository) GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Request, error) {
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

	query := fmt.Sprintf("SELECT Id, CreatorId, Info, Title, Postcode, Deadline, Status, WinningBidId FROM Requests ORDER BY deadline %s LIMIT $1 OFFSET $2;", order)

	rows, err := c.Query(ctx, query, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RequestRepo - GetAll - c.Query: %w", err)
	}
	defer rows.Close()

	var requests []entity.Request

	for rows.Next() {
		var r entity.Request
		err := rows.Scan(&r.Id, &r.CreatorId, &r.Info, &r.Title, &r.Postcode, &r.Deadline, &r.Status, &r.WinningBidId)
		if err != nil {
			return nil, fmt.Errorf("RequestRepo - GetAll - rows.Scan: %w", err)
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RequestRepo - GetAll - rows.Err: %w", err)
	}

	return &requests, nil
}

func (repo *requestRepository) CountAll(ctx context.Context) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM requests;
	`

	err = c.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("RequestRepo - CountAll - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *requestRepository) FindByCreatorId(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Request, error) {
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

	query := fmt.Sprintf("SELECT Id, CreatorId, Info, Title, Postcode, Deadline, Status, WinningBidId FROM requests WHERE CreatorId=$1 ORDER BY Deadline %s LIMIT $2 OFFSET $3;", order)

	rows, err := c.Query(ctx, query, creatorId, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RequestRepo - FindByCreatorId - c.Query: %w", err)
	}
	defer rows.Close()

	var requests []entity.Request

	for rows.Next() {
		var r entity.Request
		err := rows.Scan(&r.Id, &r.CreatorId, &r.Info, &r.Title, &r.Postcode, &r.Deadline, &r.Status, &r.WinningBidId)
		if err != nil {
			return nil, fmt.Errorf("RequestRepo - FindByCreatorId - rows.Scan: %w", err)
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RequestRepo - FindByCreatorId - rows.Err: %w", err)
	}

	return &requests, nil
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

	err = c.QueryRow(ctx, query, id).Scan(&request.Id, &request.Title, &request.Postcode, &request.Info, &request.CreatorId, &request.Deadline, &request.Status, &request.WinningBidId)
	if err != nil {
		return request, fmt.Errorf("RequestRepo - FindOneById - c.QueryRow: %w", err)
	}

	return request, nil
}

func (repo *requestRepository) CountByCreatorId(ctx context.Context, creatorId string) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM requests WHERE CreatorId=$1;
	`

	err = c.QueryRow(ctx, query, creatorId).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("RequestRepo - CountByCreatorId - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *requestRepository) UpdateWinningBidIdAndStatusById(ctx context.Context, id int, winningBidId string, status entity.RequestStatus) (int, error) {
	var requestId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return requestId, err
	}
	defer c.Release()

	const query = `
		UPDATE requests SET WinningBidId=$1, Status=$2 RETURNING Id;
	`

	err = c.QueryRow(ctx, query, winningBidId, status).Scan(&requestId)
	if err != nil {
		return requestId, fmt.Errorf("RequestRepo - UpdateWinningBidIdAndStatusById - c.QueryRow: %w", err)
	}

	return requestId, nil
}
