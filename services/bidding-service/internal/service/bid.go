package service

import (
	"context"
	"fmt"

	bidEvents "github.com/PanGan21/bidding-service/internal/events/bid"
	bidRepo "github.com/PanGan21/bidding-service/internal/repository/bid"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
)

type BidService interface {
	Create(ctx context.Context, creatorId string, requestId int, amount float64) (entity.Bid, error)
	FindOneById(ctx context.Context, id int) (entity.Bid, error)
	GetManyByRequestId(ctx context.Context, requestId int, pagination *pagination.Pagination) (*[]entity.Bid, error)
	GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Bid, error)
	CountOwn(ctx context.Context, creatorId string) (int, error)
}

type bidService struct {
	bidRepo   bidRepo.BidRepository
	bidEvents bidEvents.BidEvents
}

func NewBidService(br bidRepo.BidRepository, be bidEvents.BidEvents) BidService {
	return &bidService{bidRepo: br, bidEvents: be}
}

func (s *bidService) Create(ctx context.Context, creatorId string, requestId int, amount float64) (entity.Bid, error) {
	var newBid entity.Bid

	bidId, err := s.bidRepo.Create(ctx, creatorId, requestId, amount)
	if err != nil {
		return newBid, fmt.Errorf("BidService - Create - s.bidRepo.Create: %w", err)
	}

	newBid, err = s.bidRepo.FindOneById(ctx, bidId)
	if err != nil {
		return newBid, fmt.Errorf("BidService - Create - s.bidRepo.FindOneById: %w", err)
	}

	err = s.bidEvents.PublishBidCreated(&newBid)
	if err != nil {
		return newBid, fmt.Errorf("BidService - Create - s.bidEvents.PublishBidCreated: %w", err)
	}

	return newBid, nil
}

func (s *bidService) FindOneById(ctx context.Context, id int) (entity.Bid, error) {
	var bid entity.Bid

	bid, err := s.bidRepo.FindOneById(ctx, id)
	if err != nil {
		return bid, fmt.Errorf("BidService - FindOneById - s.bidRepo.FindOneById: %w", err)
	}

	return bid, nil
}

func (s *bidService) GetManyByRequestId(ctx context.Context, requestId int, pagination *pagination.Pagination) (*[]entity.Bid, error) {
	bids, err := s.bidRepo.FindByRequestId(ctx, requestId, pagination)
	if err != nil {
		return nil, fmt.Errorf("BidService - GetManyByRequestId - s.bidRepo.FindByRequestId: %w", err)
	}

	return bids, nil
}

func (s *bidService) GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Bid, error) {
	bids, err := s.bidRepo.FindByCreatorId(ctx, creatorId, pagination)
	if err != nil {
		return nil, fmt.Errorf("BidService - GetOwn - s.bidRepo.FindByCreatorId: %w", err)
	}

	return bids, nil
}

func (s *bidService) CountOwn(ctx context.Context, creatorId string) (int, error) {
	count, err := s.bidRepo.CountByCreatorId(ctx, creatorId)
	if err != nil {
		return 0, fmt.Errorf("BidService - CountOwn - s.bidRepo.CountByCreatorId: %w", err)
	}

	return count, nil
}
