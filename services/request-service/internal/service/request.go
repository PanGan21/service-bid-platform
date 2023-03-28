package service

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	requestEvents "github.com/PanGan21/request-service/internal/events/request"
	requestRepo "github.com/PanGan21/request-service/internal/repository/request"
)

type RequestService interface {
	Create(ctx context.Context, creatorId, info, postcode, title string, deadline int64) (entity.Request, error)
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

	var defaultStatus = entity.NewRequest
	var defaultRejectionReason = ""

	requestId, err := s.requestRepo.Create(ctx, creatorId, info, postcode, title, deadline, defaultStatus, defaultRejectionReason)
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
