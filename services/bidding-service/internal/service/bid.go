package service

import (
	"context"
	"fmt"

	bidRepo "github.com/PanGan21/bidding-service/internal/repository/bid"
	"github.com/PanGan21/pkg/entity"
)

type BidService interface {
	Create(ctx context.Context, creatorId string, requestId int, amount float64) (entity.Bid, error)
}

type bidService struct {
	bidRepo bidRepo.BidRepository
}

func NewBidService(br bidRepo.BidRepository) BidService {
	return &bidService{bidRepo: br}
}

func (s *bidService) Create(ctx context.Context, creatorId string, requestId int, amount float64) (entity.Bid, error) {
	var newBid entity.Bid

	bidId, err := s.bidRepo.Create(ctx, creatorId, requestId, amount)
	if err != nil {
		return newBid, fmt.Errorf("BidService - Create - s.bidRepo.Create: %w", err)
	}

	fmt.Println("bidId", bidId)
	newBid, err = s.bidRepo.FindOneById(ctx, bidId)
	if err != nil {
		return newBid, fmt.Errorf("BidService - Create - s.bidRepo.FindOneById: %w", err)
	}

	return newBid, nil
}
