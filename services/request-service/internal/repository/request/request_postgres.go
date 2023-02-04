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
		UPDATE requests SET WinningBidId=$1, Status=$2 WHERE Id=$3 RETURNING Id;
	`

	err = c.QueryRow(ctx, query, winningBidId, status, id).Scan(&requestId)
	if err != nil {
		return requestId, fmt.Errorf("RequestRepo - UpdateWinningBidIdAndStatusById - c.QueryRow: %w", err)
	}

	return requestId, nil
}

func (repo *requestRepository) GetAllOpenPastTime(ctx context.Context, timestamp int64, pagination *pagination.Pagination) (*[]entity.ExtendedRequest, error) {
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
		SELECT *,
		(SELECT COUNT(*) FROM bids WHERE bids.RequestId = requests.Id) as CountBids
		FROM requests
		WHERE Status=$1 AND Deadline<$2 
		ORDER BY deadline %s LIMIT $3 OFFSET $4;`,
		order)

	rows, err := c.Query(ctx, query, entity.Open, timestamp, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RequestRepo - GetAllOpenPastTime - c.Query: %w", err)
	}
	defer rows.Close()

	var requests []entity.ExtendedRequest

	for rows.Next() {
		var r entity.ExtendedRequest
		err := rows.Scan(&r.Id, &r.CreatorId, &r.Info, &r.Title, &r.Postcode, &r.Deadline, &r.Status, &r.WinningBidId, &r.BidsCount)
		if err != nil {
			return nil, fmt.Errorf("RequestRepo - GetAllOpenPastTime - rows.Scan: %w", err)
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RequestRepo - GetAllOpenPastTime - rows.Err: %w", err)
	}

	return &requests, nil
}

func (repo *requestRepository) CountAllOpenPastTime(ctx context.Context, timestamp int64) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM requests WHERE Status=$1 AND Deadline<$2;
	`

	err = c.QueryRow(ctx, query, entity.Open, timestamp).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("RequestRepo - CountAllOpenPastTime - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *requestRepository) UpdateStatusByRequestId(ctx context.Context, status entity.RequestStatus, id int) (entity.Request, error) {
	var request entity.Request

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return request, err
	}
	defer c.Release()

	const query = `
		UPDATE requests SET Status=$1 WHERE Id=$2 RETURNING *;
	`

	err = c.QueryRow(ctx, query, status, id).Scan(&request.Id, &request.CreatorId, &request.Info, &request.Title, &request.Postcode, &request.Deadline, &request.Status, &request.WinningBidId)
	if err != nil {
		return request, fmt.Errorf("RequestRepo - UpdateStatusByRequestId - c.QueryRow: %w", err)
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
		err := rows.Scan(&r.Id, &r.CreatorId, &r.Info, &r.Title, &r.Postcode, &r.Deadline, &r.Status, &r.WinningBidId)
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

func (repo *requestRepository) GetOwnAssignedByStatuses(ctx context.Context, statuses []entity.RequestStatus, userId string, pagination *pagination.Pagination) (*[]entity.BidPopulatedRequest, error) {
	var bidPopulatedRequests []entity.BidPopulatedRequest

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
		SELECT requests.Id, requests.Title, requests.Postcode, requests.Info, requests.CreatorId, requests.Deadline, requests.Status, bids.Id AS BidId, bids.Amount AS BidAmount
		FROM bids
		JOIN requests ON requests.WinningBidId = bids.Id::varchar
		WHERE bids.CreatorId = $1
		AND requests.WinningBidId IS NOT NULL
		AND requests.Status = ANY ($2)
		ORDER BY deadline %s LIMIT $3 OFFSET $4;
	`, order)

	rows, err := c.Query(ctx, query, userId, statuses, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("RequestRepo - GetOwnAssignedByStatuses - c.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.BidPopulatedRequest
		err := rows.Scan(&r.Id, &r.CreatorId, &r.Info, &r.Title, &r.Postcode, &r.Deadline, &r.Status, &r.BidId, &r.BidAmount)
		if err != nil {
			return nil, fmt.Errorf("RequestRepo - GetOwnAssignedByStatuses - rows.Scan: %w", err)
		}
		bidPopulatedRequests = append(bidPopulatedRequests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RequestRepo - GetOwnAssignedByStatuses - rows.Err: %w", err)
	}

	return &bidPopulatedRequests, nil
}

func (repo *requestRepository) CountOwnAssignedByStatuses(ctx context.Context, statuses []entity.RequestStatus, userId string) (int, error) {
	var count = 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*)
		FROM bids
		JOIN requests ON requests.WinningBidId = bids.Id::varchar
		WHERE bids.CreatorId = $1
		AND requests.WinningBidId IS NOT NULL
		AND requests.Status = ANY ($2)
	`

	err = c.QueryRow(ctx, query, userId, statuses).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("RequestRepo - CountOwnAssignedByStatuses - c.Query: %w", err)
	}

	return count, nil
}
