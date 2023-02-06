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
	GetAllOpenPastDeadline(ctx context.Context, pagination *pagination.Pagination) (*[]entity.ExtendedRequest, error)
	CountAllOpenPastDeadline(ctx context.Context) (int, error)
	UpdateStatusByRequestId(ctx context.Context, status entity.RequestStatus, id int) (entity.Request, error)
	GetAllByStatus(ctx context.Context, status entity.RequestStatus, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountAllByStatus(ctx context.Context, status entity.RequestStatus) (int, error)
	GetOwnAssignedByStatuses(ctx context.Context, statuses []entity.RequestStatus, userId string, pagination *pagination.Pagination) (*[]entity.BidPopulatedRequest, error)
	CountOwnAssignedByStatuses(ctx context.Context, statuses []entity.RequestStatus, userId string) (int, error)
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
	return (time.Now().UTC().UnixMilli() >= request.Deadline) && (request.Status == entity.Open)
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
		return request, fmt.Errorf("RequestService - UpdateWinningBid - s.requestEvents.PublishRequestUpdated: %w", err)
	}

	return request, nil
}

func (s *requestService) GetAllOpenPastDeadline(ctx context.Context, pagination *pagination.Pagination) (*[]entity.ExtendedRequest, error) {
	now := time.Now().UTC().UnixMilli()
	requests, err := s.requestRepo.GetAllOpenPastTime(ctx, now, pagination)
	if err != nil {
		return nil, fmt.Errorf("RequestService - GetAllOpenPastDeadline - s.requestRepo.GetAllOpenPastTime: %w", err)
	}

	return requests, nil
}

func (s *requestService) CountAllOpenPastDeadline(ctx context.Context) (int, error) {
	now := time.Now().UTC().UnixMilli()
	count, err := s.requestRepo.CountAllOpenPastTime(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("RequestService - CountAllOpenPastDeadline - s.requestRepo.CountAllOpenPastTime: %w", err)
	}

	return count, nil
}

func (s *requestService) UpdateStatusByRequestId(ctx context.Context, status entity.RequestStatus, id int) (entity.Request, error) {
	request, err := s.requestRepo.UpdateStatusByRequestId(ctx, status, id)
	if err != nil {
		return request, fmt.Errorf("RequestService - UpdateStatusByRequestId - s.requestRepo.UpdateStatusByRequestId: %w", err)
	}

	err = s.requestEvents.PublishRequestUpdated(&request)
	if err != nil {
		return request, fmt.Errorf("RequestService - UpdateStatusByRequestId - s.requestEvents.PublishRequestUpdated: %w", err)
	}

	return request, nil
}

func (s *requestService) GetAllByStatus(ctx context.Context, status entity.RequestStatus, pagination *pagination.Pagination) (*[]entity.Request, error) {
	requests, err := s.requestRepo.GetAllByStatus(ctx, status, pagination)
	if err != nil {
		return nil, fmt.Errorf("RequestService - GetAllByStatus - s.requestRepo.GetAllByStatus: %w", err)
	}

	return requests, nil
}

func (s *requestService) CountAllByStatus(ctx context.Context, status entity.RequestStatus) (int, error) {
	count, err := s.requestRepo.CountAllByStatus(ctx, status)
	if err != nil {
		return count, fmt.Errorf("RequestService - CountAllByStatus - s.requestRepo.CountAllByStatus: %w", err)
	}

	return count, nil
}

func (s *requestService) GetOwnAssignedByStatuses(ctx context.Context, statuses []entity.RequestStatus, userId string, pagination *pagination.Pagination) (*[]entity.BidPopulatedRequest, error) {
	bidPopulatedRequests, err := s.requestRepo.GetOwnAssignedByStatuses(ctx, statuses, userId, pagination)
	if err != nil {
		return bidPopulatedRequests, fmt.Errorf("RequestService - GetOwnAssignedByStatuses - s.requestRepo.GetOwnAssignedByStatuses: %w", err)
	}

	return bidPopulatedRequests, nil
}

func (s *requestService) CountOwnAssignedByStatuses(ctx context.Context, statuses []entity.RequestStatus, userId string) (int, error) {
	count, err := s.requestRepo.CountOwnAssignedByStatuses(ctx, statuses, userId)
	if err != nil {
		return count, fmt.Errorf("RequestService - CountOwnAssignedByStatuses - s.requestRepo.CountOwnAssignedByStatuses: %w", err)
	}

	return count, nil
}
