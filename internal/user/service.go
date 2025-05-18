package user

import (
	"context"
	"errors"
	"fmt"

	commonpb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"
	pb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/Ostap00034/course-work-backend-user-service/util/password"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExist = errors.New("пользователь с такой электронной почтой уже существует")
	ErrBadUserId        = errors.New("неверный ид пользователя")
)

// Service описывает логику UserService.
type Service interface {
	CreateUser(ctx context.Context, email, fio, role, password string) (uuid.UUID, error)
	ValidateCredentials(ctx context.Context, email, password string) (uuid.UUID, string, error)
	GetUser(ctx context.Context, id uuid.UUID) (*pb.GetUserByIdResponse, error)
	GetAllUsers(ctx context.Context) (*pb.GetUsersResponse, error)
	ChangeUser(ctx context.Context, userId string, newUser *commonpb.UserData) (*pb.GetUserByIdResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) ChangeUser(ctx context.Context, userId string, newUser *commonpb.UserData) (*pb.GetUserByIdResponse, error) {
	u, err := s.repo.Change(ctx, userId, newUser)
	if err != nil {
		return nil, err
	}

	var user commonpb.UserData

	user.Id = u.ID.String()
	user.Email = u.Email
	user.Role = u.Role.String()
	user.Fio = u.Fio

	return &pb.GetUserByIdResponse{
		User: &user,
	}, nil
}

func (s *service) CreateUser(ctx context.Context, email, fio, role, pass string) (uuid.UUID, error) {
	hash, err := password.Hash(pass)
	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}
	return s.repo.Create(ctx, email, fio, role, hash)
}

func (s *service) GetAllUsers(ctx context.Context) (*pb.GetUsersResponse, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var usersList []*commonpb.UserData

	for _, u := range users {
		usersList = append(usersList, &commonpb.UserData{
			Id:    u.ID.String(),
			Email: u.Email,
			Role:  u.Role.String(),
			Fio:   u.Fio,
		})
	}

	return &pb.GetUsersResponse{Users: usersList}, nil
}

func (s *service) ValidateCredentials(ctx context.Context, email, pass string) (uuid.UUID, string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, "", err
	}
	if err := password.Compare(u.PasswordHash, pass); err != nil {
		return uuid.Nil, "", err
	}
	// обновим updated_at
	_ = s.repo.UpdateTimestamp(ctx, u.ID)
	return u.ID, u.Role.String(), nil
}

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*pb.GetUserByIdResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var user commonpb.UserData

	user.Id = u.ID.String()
	user.Email = u.Email
	user.Role = u.Role.String()
	user.Fio = u.Fio
	user.CreatedAt = u.CreatedAt.String()
	user.UpdatedAt = u.UpdatedAt.String()

	return &pb.GetUserByIdResponse{
		User: &user,
	}, nil
}
