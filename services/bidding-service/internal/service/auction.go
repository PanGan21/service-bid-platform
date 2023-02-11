package service

import (
	"context"
	"fmt"

	auctionRepo "github.com/PanGan21/bidding-service/internal/repository/auction"
	"github.com/PanGan21/pkg/entity"
)

type AuctionService interface {
	Create(ctx context.Context, auction entity.Auction) error
	UpdateOne(ctx context.Context, auction entity.Auction) error
	IsOpenToBidByAuctionId(ctx context.Context, auctionId int) bool
}

type auctionService struct {
	auctionRepo auctionRepo.AuctionRepository
}

func NewAuctionService(rr auctionRepo.AuctionRepository) AuctionService {
	return &auctionService{auctionRepo: rr}
}

func (s *auctionService) Create(ctx context.Context, auction entity.Auction) error {
	err := s.auctionRepo.Create(ctx, auction)
	if err != nil {
		return fmt.Errorf("AuctionService - Create - s.auctionRepo.Create: %w", err)
	}

	return nil
}

func (s *auctionService) UpdateOne(ctx context.Context, auction entity.Auction) error {
	err := s.auctionRepo.UpdateOne(ctx, auction)
	if err != nil {
		return fmt.Errorf("AuctionService - UpdateOne - s.auctionRepo.UpdateOne: %w", err)
	}

	return nil
}

func (s *auctionService) IsOpenToBidByAuctionId(ctx context.Context, auctionId int) bool {
	auction, err := s.auctionRepo.FindOneById(ctx, auctionId)
	if err != nil {
		fmt.Println("AuctionService - IsOpenToBidByAuctionId - s.auctionRepo.FindOneById: %w", err)
		return false
	}

	return auction.Status == entity.Open
}
