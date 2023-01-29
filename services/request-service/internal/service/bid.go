package service

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	bidRepo "github.com/PanGan21/request-service/internal/repository/bid"
)

type BidService interface {
	Create(ctx context.Context, bid entity.Bid) error
	FindWinningBidByRequestId(ctx context.Context, requestId string) (entity.Bid, error)
}

type bidService struct {
	bidRepo bidRepo.BidRepository
}

func NewBidService(br bidRepo.BidRepository) BidService {
	return &bidService{bidRepo: br}
}

func (s *bidService) Create(ctx context.Context, bid entity.Bid) error {
	err := s.bidRepo.Create(ctx, bid)
	if err != nil {
		return fmt.Errorf("BidService - Create - s.bidRepo.Create: %w", err)
	}

	return nil
}

func (s *bidService) FindWinningBidByRequestId(ctx context.Context, requestId string) (entity.Bid, error) {
	var winnigBid entity.Bid

	bids, err := s.bidRepo.FindManyByRequestIdWithMinAmount(ctx, requestId)
	if err != nil {
		return winnigBid, fmt.Errorf("BidService - FindWinningBidByRequestId - s.bidRepo.FindOneByRequestIdWithMinAmount: %w", err)
	}

	if len(bids) != 1 {
		return winnigBid, fmt.Errorf("BidService - FindWinningBidByRequestId - s.bidRepo.FindOneByRequestIdWithMinAmount: Winning bid can only be one")
	}

	winnigBid = bids[0]

	return winnigBid, nil
}
