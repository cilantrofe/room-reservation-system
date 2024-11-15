package grpc

import (
	hotelSvc "github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"google.golang.org/grpc"
	"os"
)

type Client struct {
	Api  hotelSvc.HotelServiceClient
	conn *grpc.ClientConn
}

func New() (*Client, error) {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = "localhost:50051" // значение по умолчанию
	}

	// Устанавливаем соединение с gRPC сервером
	conn, _ := grpc.Dial(addr, grpc.WithInsecure())
	client := hotelSvc.NewHotelServiceClient(conn)

	return &Client{Api: client, conn: conn}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
