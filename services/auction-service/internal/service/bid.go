package service

import (
	"context"
	"fmt"

	bidRepo "github.com/PanGan21/auction-service/internal/repository/bid"
	"github.com/PanGan21/pkg/entity"
)

type BidService interface {
	Create(ctx context.Context, bid entity.Bid) error
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (entity.Bid, error)
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

func (s *bidService) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (entity.Bid, error) {
	var winnigBid entity.Bid

	bids, err := s.bidRepo.FindManyByAuctionIdWithMinAmount(ctx, auctionId)
	if err != nil {
		return winnigBid, fmt.Errorf("BidService - FindWinningBidByAuctionId - s.bidRepo.FindOneByAuctionIdWithMinAmount: %w", err)
	}

	if len(bids) < 1 {
		return winnigBid, fmt.Errorf("BidService - FindWinningBidByAuctionId - s.bidRepo.FindOneByAuctionIdWithMinAmount: Winning bid can only be one")
	}

	winnigBid = bids[0]

	return winnigBid, nil
}
