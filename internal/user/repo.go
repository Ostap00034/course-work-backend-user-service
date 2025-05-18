package user

import (
	"context"
	"errors"
	"time"

	commonpb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"
	"github.com/Ostap00034/course-work-backend-user-service/ent"
	"github.com/Ostap00034/course-work-backend-user-service/ent/user"
	"github.com/google/uuid"
)

var (
	ErrEmailExists      = errors.New("пользователь с такой электронной почтой уже существует")
	ErrUserNotFound     = errors.New("такой пользователь не найден")
	ErrUsersNotFound    = errors.New("пользователи не найдены")
	ErrUserUpdateFailed = errors.New("ошибка при обновлении пользователя")
	ErrUserGetFailed    = errors.New("ошибка при получении пользователя")
	ErrInvalidRole      = errors.New("некорректная роль пользователя")
)

// Repository определяет работу с базой User.
type Repository interface {
	Create(ctx context.Context, email, fio, role, passwordHash string) (uuid.UUID, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error)
	UpdateTimestamp(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]*ent.User, error)
	Change(ctx context.Context, userId string, newUser *commonpb.UserData) (*ent.User, error)
}

type repo struct {
	client *ent.Client
}

func NewRepo(client *ent.Client) Repository {
	return &repo{client: client}
}

func (r *repo) Change(ctx context.Context, userId string, newUser *commonpb.UserData) (*ent.User, error) {
	// Конвертируем строковый ID в нужный тип (если необходимо)
	id, err := uuid.Parse(userId)
	if err != nil {
		return nil, ErrBadUserId
	}

	_, err = r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Начинаем построение операции обновления
	update := r.client.User.UpdateOneID(id)

	// Обновляем поля (адаптируйте под вашу структуру)
	if newUser.GetFio() != "" {
		update = update.SetFio(newUser.GetFio())
	}

	if newUser.GetEmail() != "" {
		update = update.SetEmail(newUser.GetEmail())
	}

	if newUser.GetRole() != "" {
		// приводим строку в тип user.Role
		// проверяем, что это действительно одна из констант
		switch newUser.Role {
		case user.RoleAdmin.String(), user.RoleMaster.String(), user.RoleClient.String():
			update = update.SetRole(user.Role(newUser.Role))
		default:
			return nil, ErrInvalidRole
		}
	}

	// Выполняем обновление
	err = update.Exec(ctx)
	if err != nil {
		return nil, ErrUserUpdateFailed
	}

	// Получаем обновленную сущность
	updatedUser, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, ErrUserGetFailed
	}

	return updatedUser, nil
}
func (r *repo) Create(ctx context.Context, email, fio, role, passwordHash string) (uuid.UUID, error) {
	u, err := r.client.User.
		Create().
		SetEmail(email).
		SetFio(fio).
		SetRole(user.Role(role)).
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

func (r *repo) GetAll(ctx context.Context) ([]*ent.User, error) {
	users, err := r.client.User.Query().All(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUsersNotFound
	}
	return users, err
}

func (r *repo) UpdateTimestamp(ctx context.Context, id uuid.UUID) error {
	_, err := r.client.User.
		UpdateOneID(id).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}
