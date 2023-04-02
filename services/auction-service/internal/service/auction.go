package service

import (
	"context"
	"fmt"
	"time"

	auctionEvents "github.com/PanGan21/auction-service/internal/events/auction"
	auctionRepo "github.com/PanGan21/auction-service/internal/repository/auction"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/pagination"
)

type AuctionService interface {
	Create(ctx context.Context, auction entity.Auction) (entity.Auction, error)
	GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountAll(ctx context.Context) (int, error)
	GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountOwn(ctx context.Context, creatorId string) (int, error)
	GetById(ctx context.Context, id int) (entity.Auction, error)
	IsAllowedToResolve(ctx context.Context, auction entity.Auction) bool
	UpdateWinningBid(ctx context.Context, auction entity.Auction, winningBidId string, winnerId string, winningAmount float64) (entity.Auction, error)
	GetAllOpenPastDeadline(ctx context.Context, pagination *pagination.Pagination) (*[]entity.ExtendedAuction, error)
	CountAllOpenPastDeadline(ctx context.Context) (int, error)
	UpdateStatusByAuctionId(ctx context.Context, status entity.AuctionStatus, id int) (entity.Auction, error)
	GetAllByStatus(ctx context.Context, status entity.AuctionStatus, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountAllByStatus(ctx context.Context, status entity.AuctionStatus) (int, error)
	GetOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string, pagination *pagination.Pagination) (*[]entity.Auction, error)
	CountOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string) (int, error)
}

type auctionService struct {
	auctionRepo   auctionRepo.AuctionRepository
	auctionEvents auctionEvents.AuctionEvents
}

func NewAuctionService(rr auctionRepo.AuctionRepository, re auctionEvents.AuctionEvents) AuctionService {
	return &auctionService{auctionRepo: rr, auctionEvents: re}
}

func (s *auctionService) Create(ctx context.Context, auction entity.Auction) (entity.Auction, error) {
	var newAuction entity.Auction

	auctionId, err := s.auctionRepo.Create(ctx, auction)
	if err != nil {
		return newAuction, fmt.Errorf("AuctionService - Create - s.auctionRepo.Create: %w", err)
	}

	newAuction, err = s.auctionRepo.FindOneById(ctx, auctionId)
	if err != nil {
		return newAuction, fmt.Errorf("AuctionService - Create - s.auctionRepo.FindOneById: %w", err)
	}

	return newAuction, nil
}

func (s *auctionService) GetAll(ctx context.Context, pagination *pagination.Pagination) (*[]entity.Auction, error) {
	auctions, err := s.auctionRepo.GetAll(ctx, pagination)
	if err != nil {
		return nil, fmt.Errorf("AuctionService - GetAll - s.auctionRepo.GetAll: %w", err)
	}

	return auctions, nil
}

func (s *auctionService) CountAll(ctx context.Context) (int, error) {
	count, err := s.auctionRepo.CountAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("AuctionService - CountAll - s.auctionRepo.CountAll: %w", err)
	}

	return count, nil
}

func (s *auctionService) GetOwn(ctx context.Context, creatorId string, pagination *pagination.Pagination) (*[]entity.Auction, error) {
	auctions, err := s.auctionRepo.FindByCreatorId(ctx, creatorId, pagination)
	if err != nil {
		return nil, fmt.Errorf("AuctionService - GetOwn - s.auctionRepo.FindByCreatorId: %w", err)
	}

	return auctions, nil
}

func (s *auctionService) CountOwn(ctx context.Context, creatorId string) (int, error) {
	count, err := s.auctionRepo.CountByCreatorId(ctx, creatorId)
	if err != nil {
		return 0, fmt.Errorf("AuctionService - CountOwn - s.auctionRepo.CountByCreatorId: %w", err)
	}

	return count, nil
}

func (s *auctionService) GetById(ctx context.Context, id int) (entity.Auction, error) {
	auction, err := s.auctionRepo.FindOneById(ctx, id)
	if err != nil {
		return auction, fmt.Errorf("AuctionService - GetById - s.auctionRepo.FindOneById: %w", err)
	}

	return auction, nil
}

func (s *auctionService) IsAllowedToResolve(ctx context.Context, auction entity.Auction) bool {
	return (time.Now().UTC().UnixMilli() >= auction.Deadline) && (auction.Status == entity.Open)
}

func (s *auctionService) UpdateWinningBid(ctx context.Context, auction entity.Auction, winningBidId string, winnerId string, winningAmount float64) (entity.Auction, error) {
	if winningBidId == "" {
		return auction, fmt.Errorf("AuctionService - UpdateWinningBid: winningBid cannot be empty")
	}

	auction.WinningBidId = winningBidId
	auction.Status = entity.Assigned
	auction.WinnerId = winnerId
	auction.WinningAmount = winningAmount

	updatedAuction, err := s.auctionRepo.UpdateWinningBidIdAndStatusById(ctx, auction.Id, winningBidId, auction.Status, winnerId, winningAmount)
	if err != nil {
		return auction, fmt.Errorf("AuctionService - UpdateWinningBid - s.auctionRepo.UpdateWinningBidIdAndStatusById: %w", err)
	}

	err = s.auctionEvents.PublishAuctionUpdated(&updatedAuction)
	if err != nil {
		return updatedAuction, fmt.Errorf("AuctionService - UpdateWinningBid - s.auctionEvents.PublishAuctionUpdated: %w", err)
	}

	return updatedAuction, nil
}

func (s *auctionService) GetAllOpenPastDeadline(ctx context.Context, pagination *pagination.Pagination) (*[]entity.ExtendedAuction, error) {
	now := time.Now().UTC().UnixMilli()
	auctions, err := s.auctionRepo.GetAllOpenPastTime(ctx, now, pagination)
	if err != nil {
		return nil, fmt.Errorf("AuctionService - GetAllOpenPastDeadline - s.auctionRepo.GetAllOpenPastTime: %w", err)
	}

	return auctions, nil
}

func (s *auctionService) CountAllOpenPastDeadline(ctx context.Context) (int, error) {
	now := time.Now().UTC().UnixMilli()
	count, err := s.auctionRepo.CountAllOpenPastTime(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("AuctionService - CountAllOpenPastDeadline - s.auctionRepo.CountAllOpenPastTime: %w", err)
	}

	return count, nil
}

func (s *auctionService) UpdateStatusByAuctionId(ctx context.Context, status entity.AuctionStatus, id int) (entity.Auction, error) {
	auction, err := s.auctionRepo.UpdateStatusByAuctionId(ctx, status, id)
	if err != nil {
		return auction, fmt.Errorf("AuctionService - UpdateStatusByAuctionId - s.auctionRepo.UpdateStatusByAuctionId: %w", err)
	}

	err = s.auctionEvents.PublishAuctionUpdated(&auction)
	if err != nil {
		return auction, fmt.Errorf("AuctionService - UpdateStatusByAuctionId - s.auctionEvents.PublishAuctionUpdated: %w", err)
	}

	return auction, nil
}

func (s *auctionService) GetAllByStatus(ctx context.Context, status entity.AuctionStatus, pagination *pagination.Pagination) (*[]entity.Auction, error) {
	auctions, err := s.auctionRepo.GetAllByStatus(ctx, status, pagination)
	if err != nil {
		return nil, fmt.Errorf("AuctionService - GetAllByStatus - s.auctionRepo.GetAllByStatus: %w", err)
	}

	return auctions, nil
}

func (s *auctionService) CountAllByStatus(ctx context.Context, status entity.AuctionStatus) (int, error) {
	count, err := s.auctionRepo.CountAllByStatus(ctx, status)
	if err != nil {
		return count, fmt.Errorf("AuctionService - CountAllByStatus - s.auctionRepo.CountAllByStatus: %w", err)
	}

	return count, nil
}

func (s *auctionService) GetOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string, pagination *pagination.Pagination) (*[]entity.Auction, error) {
	auctions, err := s.auctionRepo.GetOwnAssignedByStatuses(ctx, statuses, userId, pagination)
	if err != nil {
		return auctions, fmt.Errorf("AuctionService - GetOwnAssignedByStatuses - s.auctionRepo.GetOwnAssignedByStatuses: %w", err)
	}

	return auctions, nil
}

func (s *auctionService) CountOwnAssignedByStatuses(ctx context.Context, statuses []entity.AuctionStatus, userId string) (int, error) {
	count, err := s.auctionRepo.CountOwnAssignedByStatuses(ctx, statuses, userId)
	if err != nil {
		return count, fmt.Errorf("AuctionService - CountOwnAssignedByStatuses - s.auctionRepo.CountOwnAssignedByStatuses: %w", err)
	}

	return count, nil
}
