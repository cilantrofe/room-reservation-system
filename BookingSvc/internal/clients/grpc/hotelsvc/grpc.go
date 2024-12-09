package grpc

import (
	"fmt"
	hotelSvc "github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"google.golang.org/grpc"
)

type HotelSvcClient struct {
	Api  hotelSvc.HotelServiceClient
	conn *grpc.ClientConn
}

func NewHotelClient(grpcHost, grpcPort string) (*HotelSvcClient, error) {
	address := fmt.Sprintf("%s:%s", grpcHost, grpcPort)
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Добавить ретраи мб сервис упадет??
	if err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}
	client := hotelSvc.NewHotelServiceClient(conn)
	return &HotelSvcClient{Api: client, conn: conn}, nil
}

func (c *HotelSvcClient) Close() {
	err := c.conn.Close()
	if err != nil {
		fmt.Errorf("could not close connection: %w", err)
	}
}
