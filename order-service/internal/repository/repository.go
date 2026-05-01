package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/mfuadfakhruzzaki/grpc-ecommerce/order-service/internal/repository/db"
)

type Repository struct {
	q *db.Queries
}

func New(sqlDB *sql.DB) *Repository {
	return &Repository{q: db.New(sqlDB)}
}

func (r *Repository) CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	return r.q.CreateOrder(ctx, arg)
}

func (r *Repository) GetOrderByID(ctx context.Context, id uuid.UUID) (db.Order, error) {
	return r.q.GetOrderByID(ctx, id)
}

func (r *Repository) ListOrdersByUser(ctx context.Context, arg db.ListOrdersByUserParams) ([]db.Order, error) {
	return r.q.ListOrdersByUser(ctx, arg)
}

func (r *Repository) CountOrdersByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.q.CountOrdersByUser(ctx, userID)
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, arg db.UpdateOrderStatusParams) (db.Order, error) {
	return r.q.UpdateOrderStatus(ctx, arg)
}

func (r *Repository) CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error) {
	return r.q.CreateOrderItem(ctx, arg)
}

func (r *Repository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]db.OrderItem, error) {
	return r.q.GetOrderItems(ctx, uuid.NullUUID{UUID: orderID, Valid: true})
}