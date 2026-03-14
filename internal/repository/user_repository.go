package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	db "github.com/octguy/auth-sqlc/db/sqlc"
	"github.com/octguy/auth-sqlc/internal/model"
)

// Sentinel errors — the service layer checks these, not raw pgx errors.
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrEmailDuplicate = errors.New("email already exists")
)

// UserRepository is the interface the service depends on.
// To swap backends (e.g. for testing), implement this interface.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type userRepo struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepo{q: q}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		if isDuplicate(err) {
			return ErrEmailDuplicate
		}
		return err
	}

	// write back DB-assigned fields
	user.ID = row.ID
	user.CreatedAt = row.CreatedAt
	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return toModel(row), nil
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row, err := r.q.GetUserByID(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return toModel(row), nil
}

// toModel converts a sqlc User to the domain model.
// This keeps db/sqlc isolated inside the repository layer.
func toModel(u db.User) *model.User {
	return &model.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	return strings.Contains(msg, "23505") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint")
}
