package user

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/Ostap00034/course-work-backend-user-service/util/password"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExist = errors.New("пользователь с такой электронной почтой уже существует")
)

// Service описывает логику UserService.
type Service interface {
	CreateUser(ctx context.Context, email, password string) (uuid.UUID, error)
	ValidateCredentials(ctx context.Context, email, password string) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (*pb.GetUserByIdResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateUser(ctx context.Context, email, pass string) (uuid.UUID, error) {
	hash, err := password.Hash(pass)
	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}
	return s.repo.Create(ctx, email, hash)
}

func (s *service) ValidateCredentials(ctx context.Context, email, pass string) (uuid.UUID, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, err
	}
	if err := password.Compare(u.PasswordHash, pass); err != nil {
		return uuid.Nil, err
	}
	// обновим updated_at
	_ = s.repo.UpdateTimestamp(ctx, u.ID)
	return u.ID, nil
}

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*pb.GetUserByIdResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByIdResponse{
		UserId:    u.ID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}, nil
}
