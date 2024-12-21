package service

import (
	"context"
)

type OwnerRepository interface {
	GetOwnerIdByHotelId(ctx context.Context, hotelID int) (int, error)
}

type OwnerService struct {
	ownerRepo OwnerRepository
}

func NewOwnerService(ownerRepo OwnerRepository) *OwnerService {
	return &OwnerService{ownerRepo: ownerRepo}
}

func (s *OwnerService) GetOwnerIdByHotelId(ctx context.Context, hotelID int) (int, error) {
	ownerID, err := s.ownerRepo.GetOwnerIdByHotelId(ctx, hotelID)
	if err != nil {
		return 0, err
	}
	return ownerID, err
}
