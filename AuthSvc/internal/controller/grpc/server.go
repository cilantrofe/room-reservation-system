package grpc

import (
	"context"
	"errors"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/controller"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	authSvc controller.AuthService
	Addr    string
}

func NewServer(authSvc controller.AuthService, addr string) *Server {
	return &Server{authSvc: authSvc, addr: addr}
}

func (s *Server) GetHotelierInformation(ctx context.Context, req *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error) {
	response, err := s.authSvc.GetHotelierInformation(ctx, req)
	if err != nil {
		if errors.Is(err, myerror.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, myerror.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return response, nil
}
