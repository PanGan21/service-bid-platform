package service

import (
	"context"
	"fmt"

	requestRepo "github.com/PanGan21/bidding-service/internal/repository/request"
	"github.com/PanGan21/pkg/entity"
)

type RequestService interface {
	Create(ctx context.Context, request entity.Request) error
	UpdateOne(ctx context.Context, request entity.Request) error
	IsOpenToBidByRequestId(ctx context.Context, requestId int) bool
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

func (s *requestService) UpdateOne(ctx context.Context, request entity.Request) error {
	err := s.requestRepo.UpdateOne(ctx, request)
	if err != nil {
		return fmt.Errorf("RequestService - UpdateOne - s.requestRepo.UpdateOne: %w", err)
	}

	return nil
}

func (s *requestService) IsOpenToBidByRequestId(ctx context.Context, requestId int) bool {
	request, err := s.requestRepo.FindOneById(ctx, requestId)
	if err != nil {
		fmt.Println("RequestService - IsOpenToBidByRequestId - s.requestRepo.FindOneById: %w", err)
		return false
	}

	return request.Status == entity.Open
}
