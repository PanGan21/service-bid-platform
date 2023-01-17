package service

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
	requestEvents "github.com/PanGan21/request-service/internal/events/request"
	requestRepo "github.com/PanGan21/request-service/internal/repository/request"
)

type RequestService interface {
	Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64) (entity.Request, error)
	GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Request, error)
	GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountOwn(ctx context.Context, creatorId string) (int, error)
}

type requestService struct {
	requestRepo   requestRepo.RequestRepository
	requestEvents requestEvents.RequestEvents
}

func NewRequestService(rr requestRepo.RequestRepository, re requestEvents.RequestEvents) RequestService {
	return &requestService{requestRepo: rr, requestEvents: re}
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

	err = s.requestEvents.PublishRequestCreated("request-created", &newRequest)
	if err != nil {
		return newRequest, fmt.Errorf("RequestService - Create - s.requestEvents.PublishRequestCreated: %w", err)
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

func (s *requestService) CountOwn(ctx context.Context, creatorId string) (int, error) {
	count, err := s.requestRepo.CountByCreatorId(ctx, creatorId)
	if err != nil {
		return 0, fmt.Errorf("RequestService - CountOwn - s.requestRepo.CountByCreatorId: %w", err)
	}

	return count, nil
}
