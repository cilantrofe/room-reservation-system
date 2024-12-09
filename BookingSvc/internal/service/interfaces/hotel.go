package interfaces

import (
	"context"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
)

//go:generate mockgen -source=hotel.go -destination=mocks/hotel_mock.go -package=mocks
type HotelClient interface {
	GetRoomsByHotelId(ctx context.Context, req *hotelpb.GetRoomsRequest) (*hotelpb.GetRoomsResponse, error)
}
