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
	Create(ctx context.Context, creatorId, info, postcode, title string) (entity.Request, error)
	RejectRequest(ctx context.Context, rejectionReason string, id int) (entity.Request, error)
	GetAllByStatus(ctx context.Context, status entity.RequestStatus, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountAllByStatus(ctx context.Context, status entity.RequestStatus) (int, error)
	GetManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string, pagination *pagination.Pagination) (*[]entity.Request, error)
	CountManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string) (int, error)
	ApproveRequestById(ctx context.Context, id int) (entity.Request, error)
}

type requestService struct {
	requestRepo   requestRepo.RequestRepository
	requestEvents requestEvents.RequestEvents
}

func NewRequestService(rr requestRepo.RequestRepository, re requestEvents.RequestEvents) RequestService {
	return &requestService{requestRepo: rr, requestEvents: re}
}

func (s *requestService) Create(ctx context.Context, creatorId, info, postcode, title string) (entity.Request, error) {
	var newRequest entity.Request

	var defaultStatus = entity.NewRequest
	var defaultRejectionReason = ""

	requestId, err := s.requestRepo.Create(ctx, creatorId, info, postcode, title, defaultStatus, defaultRejectionReason)
	if err != nil {
		return newRequest, fmt.Errorf("RequestService - Create - s.requestRepo.Create: %w", err)
	}

	newRequest, err = s.requestRepo.FindOneById(ctx, requestId)
	if err != nil {
		return newRequest, fmt.Errorf("RequestService - Create - s.requestRepo.FindOneById: %w", err)
	}

	return newRequest, nil
}

func (s *requestService) RejectRequest(ctx context.Context, rejectionReason string, id int) (entity.Request, error) {
	request, err := s.requestRepo.UpdateStatusAndRejectionReasonById(ctx, id, entity.RejectedRequest, rejectionReason)
	if err != nil {
		return request, fmt.Errorf("RequestService - RejectRequest - s.requestRepo.UpdateStatusAndRejectionReason: %w", err)
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

func (s *requestService) GetManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string, pagination *pagination.Pagination) (*[]entity.Request, error) {
	requests, err := s.requestRepo.GetManyByStatusByUserId(ctx, status, userId, pagination)
	if err != nil {
		return nil, fmt.Errorf("RequestService - GetOwnByStatus - s.requestRepo.GetOwnByStatus: %w", err)
	}

	return requests, nil
}

func (s *requestService) CountManyByStatusByUserId(ctx context.Context, status entity.RequestStatus, userId string) (int, error) {
	count, err := s.requestRepo.CountManyByStatusByUserId(ctx, status, userId)
	if err != nil {
		return count, fmt.Errorf("RequestService - CountManyByStatusByUserId - s.requestRepo.CountManyByStatusByUserId: %w", err)
	}

	return count, nil
}

func (s *requestService) ApproveRequestById(ctx context.Context, id int) (entity.Request, error) {
	request, err := s.requestRepo.UpdateStatusById(ctx, id, entity.ApprovedRequest)
	if err != nil {
		return request, fmt.Errorf("RequestService - ApproveRequestById - s.requestRepo.UpdateStatusById: %w", err)
	}

	err = s.requestEvents.PublishRequestApproved(&request)
	if err != nil {
		return request, fmt.Errorf("RequestService - Create - s.requestEvents.PublishRequestApproved: %w", err)
	}

	return request, nil
}
