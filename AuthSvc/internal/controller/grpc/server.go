package grpc

import (
	"context"
	"errors"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/controller"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	authSvc controller.AuthService
	Addr    string
}

func NewServer(authSvc controller.AuthService, addr string) *Server {
	return &Server{authSvc: authSvc, Addr: addr}
}

func (s *Server) GetHotelierInformation(ctx context.Context, req *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.Int("request.chat_id", int(req.GetOwnerID())),
	)

	response, err := s.authSvc.GetHotelierInformation(ctx, req)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, myerror.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, myerror.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return response, nil
}
