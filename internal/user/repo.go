package user

import (
	"context"
	"errors"
	"time"

	"github.com/Ostap00034/course-work-backend-user-service/ent"
	"github.com/Ostap00034/course-work-backend-user-service/ent/user"
	"github.com/google/uuid"
)

var (
	ErrEmailExists  = errors.New("пользователь с такой электронной почтой уже существует")
	ErrUserNotFound = errors.New("такой пользователь не найден")
)

// Repository определяет работу с базой User.
type Repository interface {
	Create(ctx context.Context, email, passwordHash string) (uuid.UUID, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error)
	UpdateTimestamp(ctx context.Context, id uuid.UUID) error
}

type repo struct {
	client *ent.Client
}

func NewRepo(client *ent.Client) Repository {
	return &repo{client: client}
}

func (r *repo) Create(ctx context.Context, email, passwordHash string) (uuid.UUID, error) {
	u, err := r.client.User.
		Create().
		SetEmail(email).
		SetPasswordHash(passwordHash).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return uuid.Nil, ErrEmailExists
		}
		return uuid.Nil, err
	}
	return u.ID, nil
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	return u, err
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	u, err := r.client.User.
		Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	return u, err
}

func (r *repo) UpdateTimestamp(ctx context.Context, id uuid.UUID) error {
	_, err := r.client.User.
		UpdateOneID(id).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}
