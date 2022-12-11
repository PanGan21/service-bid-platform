package service

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
	"github.com/PanGan21/request-service/internal/repository/request"
)

type RequestService interface {
	Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64) (entity.Request, error)
	GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Request, error)
	GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Request, error)
}

type requestService struct {
	requestRepo request.RequestRepository
}

func NewRequestService(rr request.RequestRepository) RequestService {
	return &requestService{requestRepo: rr}
}

func (s *requestService) Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64) (entity.Request, error) {
	var newRequest entity.Request

	requestId, err := s.requestRepo.Create(ctx, creatorId, info, postcode, title, deadline)
	if err != nil {
		return newRequest, fmt.Errorf("RequestService - Create - s.requestRepo.Create: %w", err)
	}

	newRequest, err = s.requestRepo.FindOneById(ctx, requestId)
	if err != nil {
		return newRequest, fmt.Errorf("RequestService - Create - s.requestRepo.FindOneById: %w", err)
	}

	return newRequest, nil
}

func (s *requestService) GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Request, error) {
	requests, err := s.requestRepo.GetAll(ctx, pagination)
	if err != nil {
		return nil, fmt.Errorf("RequestService - GetAll - s.requestRepo.GetAll: %w", err)
	}

	return requests, nil
}

func (s *requestService) GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Request, error) {
	requests, err := s.requestRepo.FindByCreatorId(ctx, creatorId, pagination)
	if err != nil {
		return nil, fmt.Errorf("RequestService - GetOwn - s.requestRepo.FindByCreatorId: %w", err)
	}

	return requests, nil
}
