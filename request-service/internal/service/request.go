package service

import (
	"context"
	"fmt"

	"github.com/PanGan21/request-service/internal/entity"
	"github.com/PanGan21/request-service/internal/repository/request"
	"github.com/google/uuid"
)

type RequestService interface {
	Create(ctx context.Context, creatorId string, title string, postcode string, info string, deadline int64) (*entity.Request, error)
}

type requestService struct {
	requestRepo request.RequestRepository
}

func NewRequestService(rr request.RequestRepository) RequestService {
	return &requestService{requestRepo: rr}
}

func (s *requestService) Create(ctx context.Context, creatorId string, title string, postcode string, info string, deadline int64) (*entity.Request, error) {
	id := uuid.New()

	request := &entity.Request{
		Id:        id,
		Title:     title,
		Postcode:  postcode,
		Info:      info,
		CreatorId: creatorId,
		Deadline:  deadline,
	}
	err := s.requestRepo.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("RequestService - Create - s.requestRepo.Create: %w", err)
	}

	return request, nil
}
