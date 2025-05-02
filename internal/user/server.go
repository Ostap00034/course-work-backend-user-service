package user

import (
	"context"

	userpb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Встраиваем UnimplementedUserServiceServer
type Server struct {
	userpb.UnimplementedUserServiceServer
	svc Service
}

func NewServer(s Service) *Server {
	return &Server{svc: s}
}

func (s *Server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	id, err := s.svc.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userpb.CreateUserResponse{UserId: id.String()}, nil
}

func (s *Server) ValidateCredentials(ctx context.Context, req *userpb.ValidateCredentialsRequest) (*userpb.ValidateCredentialsResponse, error) {
	id, err := s.svc.ValidateCredentials(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &userpb.ValidateCredentialsResponse{UserId: id.String()}, nil
}

func (s *Server) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserByIdResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad userId")
	}
	return s.svc.GetUser(ctx, uid)
}
