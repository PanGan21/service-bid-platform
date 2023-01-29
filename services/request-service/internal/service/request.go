package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
	requestEvents "github.com/PanGan21/request-service/internal/events/request"
	requestRepo "github.com/PanGan21/request-service/internal/repository/request"
)

type RequestService interface {
	Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64) (entity.Request, error)
	GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountAll(ctx context.Context) (int, error)
	GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountOwn(ctx context.Context, creatorId string) (int, error)
	GetById(ctx context.Context, id int) (entity.Request, error)
	IsAllowedToResolve(ctx context.Context, request entity.Request) bool
	UpdateWinningBid(ctx context.Context, request entity.Request, winningBidId string) (entity.Request, error)
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

	err = s.requestEvents.PublishRequestCreated(&newRequest)
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

func (s *requestService) CountAll(ctx context.Context) (int, error) {
	count, err := s.requestRepo.CountAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("RequestService - CountAll - s.requestRepo.CountAll: %w", err)
	}

	return count, nil
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

func (s *requestService) GetById(ctx context.Context, id int) (entity.Request, error) {
	request, err := s.requestRepo.FindOneById(ctx, id)
	if err != nil {
		return request, fmt.Errorf("RequestService - GetById - s.requestRepo.FindOneById: %w", err)
	}

	return request, nil
}

func (s *requestService) IsAllowedToResolve(ctx context.Context, request entity.Request) bool {
	return (time.Now().Unix() >= request.Deadline) && (request.Status == entity.Open)
}

func (s *requestService) UpdateWinningBid(ctx context.Context, request entity.Request, winningBidId string) (entity.Request, error) {
	if winningBidId == "" {
		return request, fmt.Errorf("RequestService - UpdateWinningBid: winningBid cannot be empty")
	}

	request.WinningBidId = winningBidId
	request.Status = entity.Assigned

	_, err := s.requestRepo.UpdateWinningBidIdAndStatusById(ctx, request.Id, winningBidId, request.Status)
	if err != nil {
		return request, fmt.Errorf("RequestService - UpdateWinningBid - s.requestRepo.UpdateWinningBidIdAndStatusById: %w", err)
	}

	err = s.requestEvents.PublishRequestUpdated(&request)
	if err != nil {
		return request, fmt.Errorf("RequestService - UpdateWinningBid - s.requestEvents.UpdateWinningBidIdAndStatusById: %w", err)
	}

	return request, nil
}
