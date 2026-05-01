package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/repository/db"
)

type Repository struct {
	q *db.Queries
}

func New(sqlDB *sql.DB) *Repository {
	return &Repository{q: db.New(sqlDB)}
}

func (r *Repository) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	return r.q.CreateProduct(ctx, arg)
}

func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (db.Product, error) {
	return r.q.GetProductByID(ctx, id)
}

func (r *Repository) ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Product, error) {
	return r.q.ListProducts(ctx, arg)
}

func (r *Repository) CountProducts(ctx context.Context) (int64, error) {
	return r.q.CountProducts(ctx)
}

func (r *Repository) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	return r.q.UpdateProduct(ctx, arg)
}

func (r *Repository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteProduct(ctx, id)
}

func (r *Repository) GetStock(ctx context.Context, id uuid.UUID) (sql.NullInt32, error) {
	return r.q.GetStock(ctx, id)
}

func (r *Repository) DeductStock(ctx context.Context, arg db.DeductStockParams) (sql.NullInt32, error) {
	return r.q.DeductStock(ctx, arg)
}