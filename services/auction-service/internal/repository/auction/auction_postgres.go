package auction

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
	"github.com/PanGan21/pkg/postgres"
)

type auctionRepository struct {
	db postgres.Postgres
}

func NewAuctionRepository(db postgres.Postgres) *auctionRepository {
	return &auctionRepository{db: db}
}

func (repo *auctionRepository) Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64, status entity.AuctionStatus, winningBidId string, rejectionReason string) (int, error) {
	var auctionId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return auctionId, err
	}
	defer c.Release()

	const query = `
  		INSERT INTO auctions (CreatorId, Info, Postcode, Title, Deadline, Status, WinningBidId, RejectionReason) 
  		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING Id;
	`

	c.QueryRow(ctx, query, creatorId, info, postcode, title, deadline, status, winningBidId, rejectionReason).Scan(&auctionId)
	if err != nil {
		return auctionId, fmt.Errorf("AuctionRepo - Create - c.QueryRow: %w", err)
	}

	return auctionId, nil
}

func (repo *auctionRepository) GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Auction, error) {
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

	query := fmt.Sprintf("SELECT * FROM Auctions ORDER BY deadline %s LIMIT $1 OFFSET $2;", order)

	rows, err := c.Query(ctx, query, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetAll - c.Query: %w", err)
	}
	defer rows.Close()

	var auctions []entity.Auction

	for rows.Next() {
		var r entity.Auction
		err := rows.Scan(&r.Id, &r.Title, &r.Postcode, &r.Info, &r.CreatorId, &r.Deadline, &r.Status, &r.WinningBidId, &r.RejectionReason)
		if err != nil {
			return nil, fmt.Errorf("AuctionRepo - GetAll - rows.Scan: %w", err)
		}
		auctions = append(auctions, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetAll - rows.Err: %w", err)
	}

	return &auctions, nil
}

func (repo *auctionRepository) CountAll(ctx context.Context) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM auctions;
	`

	err = c.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("AuctionRepo - CountAll - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *auctionRepository) FindByCreatorId(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Auction, error) {
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

	query := fmt.Sprintf("SELECT * FROM auctions WHERE CreatorId=$1 ORDER BY Deadline %s LIMIT $2 OFFSET $3;", order)

	rows, err := c.Query(ctx, query, creatorId, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("AuctionRepo - FindByCreatorId - c.Query: %w", err)
	}
	defer rows.Close()

	var auctions []entity.Auction

	for rows.Next() {
		var r entity.Auction
		err := rows.Scan(&r.Id, &r.Title, &r.Postcode, &r.Info, &r.CreatorId, &r.Deadline, &r.Status, &r.WinningBidId, &r.RejectionReason)
		if err != nil {
			return nil, fmt.Errorf("AuctionRepo - FindByCreatorId - rows.Scan: %w", err)
		}
		auctions = append(auctions, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AuctionRepo - FindByCreatorId - rows.Err: %w", err)
	}

	return &auctions, nil
}

func (repo *auctionRepository) FindOneById(ctx context.Context, id int) (entity.Auction, error) {
	var auction entity.Auction

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return auction, err
	}
	defer c.Release()

	const query = `
		SELECT * FROM auctions WHERE Id=$1;
	`

	err = c.QueryRow(ctx, query, id).Scan(&auction.Id, &auction.Title, &auction.Postcode, &auction.Info, &auction.CreatorId, &auction.Deadline, &auction.Status, &auction.WinningBidId, &auction.RejectionReason)
	if err != nil {
		return auction, fmt.Errorf("AuctionRepo - FindOneById - c.QueryRow: %w", err)
	}

	return auction, nil
}

func (repo *auctionRepository) CountByCreatorId(ctx context.Context, creatorId string) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM auctions WHERE CreatorId=$1;
	`

	err = c.QueryRow(ctx, query, creatorId).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("AuctionRepo - CountByCreatorId - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *auctionRepository) UpdateWinningBidIdAndStatusById(ctx context.Context, id int, winningBidId string, status entity.AuctionStatus) (int, error) {
	var auctionId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return auctionId, err
	}
	defer c.Release()

	const query = `
		UPDATE auctions SET WinningBidId=$1, Status=$2 WHERE Id=$3 RETURNING Id;
	`

	err = c.QueryRow(ctx, query, winningBidId, status, id).Scan(&auctionId)
	if err != nil {
		return auctionId, fmt.Errorf("AuctionRepo - UpdateWinningBidIdAndStatusById - c.QueryRow: %w", err)
	}

	return auctionId, nil
}

func (repo *auctionRepository) GetAllOpenPastTime(ctx context.Context, timestamp int64, pagination *pagination.Pagination) (*[]entity.ExtendedAuction, error) {
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
		SELECT Id, CreatorId, Info, Title, Postcode, Deadline, Status, WinningBidId,
		(SELECT COUNT(*) FROM bids WHERE bids.AuctionId = auctions.Id) as CountBids
		FROM auctions
		WHERE Status=$1 AND Deadline<$2 
		ORDER BY deadline %s LIMIT $3 OFFSET $4;`,
		order)

	rows, err := c.Query(ctx, query, entity.Open, timestamp, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetAllOpenPastTime - c.Query: %w", err)
	}
	defer rows.Close()

	var auctions []entity.ExtendedAuction

	for rows.Next() {
		var r entity.ExtendedAuction
		err := rows.Scan(&r.Id, &r.CreatorId, &r.Info, &r.Title, &r.Postcode, &r.Deadline, &r.Status, &r.WinningBidId, &r.BidsCount)
		if err != nil {
			return nil, fmt.Errorf("AuctionRepo - GetAllOpenPastTime - rows.Scan: %w", err)
		}
		auctions = append(auctions, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetAllOpenPastTime - rows.Err: %w", err)
	}

	return &auctions, nil
}

func (repo *auctionRepository) CountAllOpenPastTime(ctx context.Context, timestamp int64) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM auctions WHERE Status=$1 AND Deadline<$2;
	`

	err = c.QueryRow(ctx, query, entity.Open, timestamp).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("AuctionRepo - CountAllOpenPastTime - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *auctionRepository) UpdateStatusByAuctionId(ctx context.Context, status entity.AuctionStatus, id int) (entity.Auction, error) {
	var auction entity.Auction

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return auction, err
	}
	defer c.Release()

	const query = `
		UPDATE auctions SET Status=$1 WHERE Id=$2 RETURNING *;
	`

	err = c.QueryRow(ctx, query, status, id).Scan(&auction.Id, &auction.CreatorId, &auction.Info, &auction.Title, &auction.Postcode, &auction.Deadline, &auction.Status, &auction.WinningBidId, &auction.RejectionReason)
	if err != nil {
		return auction, fmt.Errorf("AuctionRepo - UpdateStatusByAuctionId - c.QueryRow: %w", err)
	}

	return auction, nil
}

func (repo *auctionRepository) GetAllByStatus(ctx context.Context, status entity.AuctionStatus, pagination *pagination.Pagination) (*[]entity.Auction, error) {
	var auctions []entity.Auction

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

	query := fmt.Sprintf("SELECT * FROM auctions WHERE Status=$1 ORDER BY Deadline %s LIMIT $2 OFFSET $3;", order)

	rows, err := c.Query(ctx, query, status, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetAllByStatus - c.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.Auction
		err := rows.Scan(&r.Id, &r.Title, &r.Postcode, &r.Info, &r.CreatorId, &r.Deadline, &r.Status, &r.WinningBidId, &r.RejectionReason)
		if err != nil {
			return nil, fmt.Errorf("AuctionRepo - GetAllByStatus - rows.Scan: %w", err)
		}
		auctions = append(auctions, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetAllByStatus - rows.Err: %w", err)
	}

	return &auctions, nil
}

func (repo *auctionRepository) CountAllByStatus(ctx context.Context, status entity.AuctionStatus) (int, error) {
	count := 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*) FROM auctions WHERE Status=$1;
	`

	err = c.QueryRow(ctx, query, status).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("AuctionRepo - CountAllByStatus - c.QueryRow: %w", err)
	}

	return count, nil
}

func (repo *auctionRepository) GetOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string, pagination *pagination.Pagination) (*[]entity.BidPopulatedAuction, error) {
	var bidPopulatedAuctions []entity.BidPopulatedAuction

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
		SELECT auctions.Id, auctions.Title, auctions.Postcode, auctions.Info, auctions.CreatorId, auctions.Deadline, auctions.Status, bids.Id AS BidId, bids.Amount AS BidAmount
		FROM bids
		JOIN auctions ON auctions.WinningBidId = bids.Id::varchar
		WHERE bids.CreatorId = $1
		AND auctions.WinningBidId IS NOT NULL
		AND auctions.Status = ANY ($2)
		ORDER BY deadline %s LIMIT $3 OFFSET $4;
	`, order)

	rows, err := c.Query(ctx, query, userId, statuses, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetOwnAssignedByStatuses - c.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.BidPopulatedAuction
		err := rows.Scan(&r.Id, &r.CreatorId, &r.Info, &r.Title, &r.Postcode, &r.Deadline, &r.Status, &r.BidId, &r.BidAmount)
		if err != nil {
			return nil, fmt.Errorf("AuctionRepo - GetOwnAssignedByStatuses - rows.Scan: %w", err)
		}
		bidPopulatedAuctions = append(bidPopulatedAuctions, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AuctionRepo - GetOwnAssignedByStatuses - rows.Err: %w", err)
	}

	return &bidPopulatedAuctions, nil
}

func (repo *auctionRepository) CountOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string) (int, error) {
	var count = 0

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return count, err
	}
	defer c.Release()

	const query = `
		SELECT COUNT(*)
		FROM bids
		JOIN auctions ON auctions.WinningBidId = bids.Id::varchar
		WHERE bids.CreatorId = $1
		AND auctions.WinningBidId IS NOT NULL
		AND auctions.Status = ANY ($2)
	`

	err = c.QueryRow(ctx, query, userId, statuses).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("AuctionRepo - CountOwnAssignedByStatuses - c.Query: %w", err)
	}

	return count, nil
}

func (repo *auctionRepository) UpdateStatusAndRejectionReasonById(ctx context.Context, id int, status entity.AuctionStatus, rejectionReason string) (entity.Auction, error) {
	var auction entity.Auction

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return auction, err
	}
	defer c.Release()

	const query = `
		UPDATE auctions SET Status=$1, RejectionReason=$2 WHERE Id=$3 RETURNING *;
	`

	err = c.QueryRow(ctx, query, status, rejectionReason, id).Scan(&auction.Id, &auction.Title, &auction.Postcode, &auction.Info, &auction.CreatorId, &auction.Deadline, &auction.Status, &auction.WinningBidId, &auction.RejectionReason)
	if err != nil {
		return auction, fmt.Errorf("AuctionRepo - UpdateStatusAndRejectionReasonById - c.QueryRow: %w", err)
	}

	return auction, nil
}
