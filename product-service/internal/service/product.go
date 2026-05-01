package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/google/uuid"
	db "github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/repository/db"
)

type Repo interface {
	CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (db.Product, error)
	ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Product, error)
	CountProducts(ctx context.Context) (int64, error)
	UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetStock(ctx context.Context, id uuid.UUID) (sql.NullInt32, error)
	DeductStock(ctx context.Context, arg db.DeductStockParams) (sql.NullInt32, error)
}

type ProductService struct {
	repo Repo
}

func New(repo Repo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, name, description string, price float64, categoryID string, stockQty int32) (db.Product, error) {
	priceStr := strconv.FormatFloat(price, 'f', 2, 64)
	arg := db.CreateProductParams{
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		Price:       priceStr,
		StockQty:    sql.NullInt32{Int32: stockQty, Valid: true},
	}
	if categoryID != "" {
		catUUID, err := uuid.Parse(categoryID)
		if err == nil {
			arg.CategoryID = uuid.NullUUID{UUID: catUUID, Valid: true}
		}
	}
	return s.repo.CreateProduct(ctx, arg)
}

func (s *ProductService) Get(ctx context.Context, id uuid.UUID) (db.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, page, limit int32) ([]db.Product, int64, error) {
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	products, err := s.repo.ListProducts(ctx, db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountProducts(ctx)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, name, description string, price float64) (db.Product, error) {
	priceStr := strconv.FormatFloat(price, 'f', 2, 64)
	return s.repo.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          id,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		Price:       priceStr,
	})
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) CheckStock(ctx context.Context, productID uuid.UUID, quantity int32) (bool, int32, error) {
	stock, err := s.repo.GetStock(ctx, productID)
	if err != nil {
		return false, 0, err
	}
	available := stock.Valid && stock.Int32 >= quantity
	return available, stock.Int32, nil
}

func (s *ProductService) DeductStock(ctx context.Context, productID uuid.UUID, quantity int32) (bool, int32, error) {
	remaining, err := s.repo.DeductStock(ctx, db.DeductStockParams{
		ID:       productID,
		StockQty: sql.NullInt32{Int32: quantity, Valid: true},
	})
	if err != nil {
		return false, 0, errors.New("insufficient stock or product not found")
	}
	return true, remaining.Int32, nil
}