package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/mfuadfakhruzzaki/grpc-ecommerce/user-service/internal/repository/db"
)

type Repository struct {
	q *db.Queries
}

func New(sqlDB *sql.DB) *Repository {
	return &Repository{
		q: db.New(sqlDB),
	}
}

func (r *Repository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *Repository) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	return r.q.UpdateUser(ctx, arg)
}