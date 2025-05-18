package user

import (
	"context"

	userpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Встраиваем UnimplementedUserServiceServer
type Server struct {
	userpbv1.UnimplementedUserServiceServer
	svc Service
}

func NewServer(s Service) *Server {
	return &Server{svc: s}
}

func (s *Server) CreateUser(ctx context.Context, req *userpbv1.CreateUserRequest) (*userpbv1.CreateUserResponse, error) {
	id, err := s.svc.CreateUser(ctx, req.Email, req.Fio, req.Role, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userpbv1.CreateUserResponse{UserId: id.String()}, nil
}

func (s *Server) ValidateCredentials(ctx context.Context, req *userpbv1.ValidateCredentialsRequest) (*userpbv1.ValidateCredentialsResponse, error) {
	id, role, err := s.svc.ValidateCredentials(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &userpbv1.ValidateCredentialsResponse{UserId: id.String(), Role: role}, nil
}

func (s *Server) GetUserById(ctx context.Context, req *userpbv1.GetUserByIdRequest) (*userpbv1.GetUserByIdResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrBadUserId.Error())
	}
	return s.svc.GetUser(ctx, uid)
}

func (s *Server) ChangeUser(ctx context.Context, req *userpbv1.ChangeUserRequest) (*userpbv1.GetUserByIdResponse, error) {
	return s.svc.ChangeUser(ctx, req.UserId, req.User)
}

func (s *Server) GetUsers(ctx context.Context, req *userpbv1.GetUsersRequest) (*userpbv1.GetUsersResponse, error) {
	users, err := s.svc.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return users, err
}
