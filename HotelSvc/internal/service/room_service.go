package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/models"
)

type RoomRepository interface {
	GetRoomsByHotelId(id int) ([]*hotelpb.Room, error)
	AddRoomType(roomType models.RoomType) error
	AddRoom(room models.Room) error
}

type RoomService struct {
	roomRepo RoomRepository
}

// NewRoomService создает новый экземпляр RoomService.
func NewRoomService(roomRepo RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) GetRoomsByHotelId(id int) ([]*hotelpb.Room, error) {
	rooms, err := s.roomRepo.GetRoomsByHotelId(id)
	if err != nil {
		return nil, err
	}
	return rooms, err
}

func (s *RoomService) AddRoom(ctx context.Context, room models.Room) error {
	return s.roomRepo.AddRoom(room)
}

func (s *RoomService) AddRoomType(roomType models.RoomType) error {
	return s.roomRepo.AddRoomType(roomType)
}
