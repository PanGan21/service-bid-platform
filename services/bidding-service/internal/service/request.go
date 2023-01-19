package service

import (
	"context"
	"fmt"

	requestRepo "github.com/PanGan21/bidding-service/internal/repository/request"
	"github.com/PanGan21/pkg/entity"
)

type RequestService interface {
	Create(ctx context.Context, request entity.Request) error
}

type requestService struct {
	requestRepo requestRepo.RequestRepository
}

func NewRequestService(rr requestRepo.RequestRepository) RequestService {
	return &requestService{requestRepo: rr}
}

func (s *requestService) Create(ctx context.Context, request entity.Request) error {
	err := s.requestRepo.Create(ctx, request)
	if err != nil {
		return fmt.Errorf("RequestService - Create - s.requestRepo.Create: %w", err)
	}

	return nil
}
