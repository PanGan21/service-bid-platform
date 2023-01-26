package service

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	bidRepo "github.com/PanGan21/request-service/internal/repository/bid"
)

type BidService interface {
	Create(ctx context.Context, bid entity.Bid) error
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
