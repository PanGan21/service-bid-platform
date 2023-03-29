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

func (repo *requestRepository) UpdateStatusAndRejectionReasonById(ctx context.Context, id int, status entity.RequestStatus, rejectionReason string) (entity.Request, error) {
	var request entity.Request

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return request, err
	}
	defer c.Release()

	const query = `
		UPDATE requests SET Status=$1, RejectionReason=$2 WHERE Id=$3 RETURNING *;
	`

	err = c.QueryRow(ctx, query, status, rejectionReason, id).Scan(&request.Id, &request.Title, &request.Postcode, &request.Info, &request.CreatorId, &request.Deadline, &request.Status, &request.RejectionReason)
	if err != nil {
		return request, fmt.Errorf("RequestRepo - UpdateStatusAndRejectionReasonById - c.QueryRow: %w", err)
	}

	return request, nil
}

func (repo *requestRepository) GetAllByStatus(ctx context.Context, status entity.RequestStatus, pagination *pagination.Pagination) (*[]entity.Request, error) {
	var requests []entity.Request

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

	query := fmt.Sprintf("SELECT * FROM requests WHERE Status=$1 ORDER BY Deadline %s LIMIT $2 OFFSET $3;", order)

	rows, err := c.Query(ctx, query, status, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RequestRepo - GetAllByStatus - c.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.Request
		err := rows.Scan(&r.Id, &r.Title, &r.Postcode, &r.Info, &r.CreatorId, &r.Deadline, &r.Status, &r.RejectionReason)
		if err != nil {
			return nil, fmt.Errorf("RequestRepo - GetAllByStatus - rows.Scan: %w", err)
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RequestRepo - GetAllByStatus - rows.Err: %w", err)
	}

	return &requests, nil
}

func (repo *requestRepository) CountAllByStatus(ctx context.Context, status entity.RequestStatus) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM requests WHERE Status=$1;
	`

	err = c.QueryRow(ctx, query, status).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("RequestRepo - CountAllByStatus - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *requestRepository) GetManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string, pagination *pagination.Pagination) (*[]entity.Request, error) {
	var requests []entity.Request

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

	query := fmt.Sprintf(`
		SELECT * FROM requests
		WHERE CreatorId=$1 AND Status=$2
		ORDER BY deadline %s LIMIT $3 OFFSET $4;
	`, order)

	rows, err := c.Query(ctx, query, userId, status, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RequestRepo - GetManyByStatusByUserId - c.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.Request
		err := rows.Scan(&r.Id, &r.Title, &r.Postcode, &r.Info, &r.CreatorId, &r.Deadline, &r.Status, &r.RejectionReason)
		if err != nil {
			return nil, fmt.Errorf("RequestRepo - GetManyByStatusByUserId - rows.Scan: %w", err)
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RequestRepo - GetManyByStatusByUserId - rows.Err: %w", err)
	}

	return &requests, nil
}

func (repo *requestRepository) CountManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string) (int, error) {
	var count = 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM requests
		WHERE CreatorId=$1 AND Status=$2
	`

	err = c.QueryRow(ctx, query, userId, status).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("RequestRepo - CountManyByStatusByUserId - c.Query: %w", err)
	}

	return count, nil
}
