package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	db "github.com/mfuadfakhruzzaki/grpc-ecommerce/order-service/internal/repository/db"
	pbproduct "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/product/v1"
)

type Repo interface {
	CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (db.Order, error)
	ListOrdersByUser(ctx context.Context, arg db.ListOrdersByUserParams) ([]db.Order, error)
	CountOrdersByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	UpdateOrderStatus(ctx context.Context, arg db.UpdateOrderStatusParams) (db.Order, error)
	CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error)
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]db.OrderItem, error)
}

type OrderService struct {
	repo        Repo
	productConn pbproduct.ProductServiceClient
}

func New(repo Repo, productClient pbproduct.ProductServiceClient) *OrderService {
	return &OrderService{repo: repo, productConn: productClient}
}

type OrderItemInput struct {
	ProductID string
	Quantity  int32
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, items []OrderItemInput) (db.Order, []db.OrderItem, error) {
	// Validasi dan deduct stok per item, hitung total
	var totalAmount float64
	type itemDetail struct {
		productID uuid.UUID
		quantity  int32
		price     float64
	}
	var details []itemDetail

	for _, item := range items {
		pid, err := uuid.Parse(item.ProductID)
		if err != nil {
			return db.Order{}, nil, fmt.Errorf("invalid product id: %s", item.ProductID)
		}

		// Ambil harga + cek stok dalam satu call (eliminasi CheckStock terpisah)
		productRes, err := s.productConn.GetProduct(ctx, &pbproduct.GetProductReq{Id: item.ProductID})
		if err != nil {
			return db.Order{}, nil, fmt.Errorf("product not found: %s", item.ProductID)
		}
		if productRes.Product.StockQty < item.Quantity {
			return db.Order{}, nil, fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}

		// Deduct stok (atomik di DB: WHERE stock_qty >= $2)
		_, err = s.productConn.DeductStock(ctx, &pbproduct.DeductStockReq{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
		if err != nil {
			return db.Order{}, nil, fmt.Errorf("failed to deduct stock: %v", err)
		}

		price := productRes.Product.Price
		totalAmount += price * float64(item.Quantity)
		details = append(details, itemDetail{pid, item.Quantity, price})
	}

	// Buat order
	order, err := s.repo.CreateOrder(ctx, db.CreateOrderParams{
		UserID:      userID,
		TotalAmount: strconv.FormatFloat(totalAmount, 'f', 2, 64),
	})
	if err != nil {
		return db.Order{}, nil, err
	}

	// Buat order items
	var orderItems []db.OrderItem
	for _, d := range details {
		item, err := s.repo.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:      uuid.NullUUID{UUID: order.ID, Valid: true},
			ProductID:    d.productID,
			Quantity:     d.quantity,
			PriceAtOrder: strconv.FormatFloat(d.price, 'f', 2, 64),
		})
		if err != nil {
			return db.Order{}, nil, err
		}
		orderItems = append(orderItems, item)
	}

	return order, orderItems, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (db.Order, []db.OrderItem, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return db.Order{}, nil, err
	}
	items, err := s.repo.GetOrderItems(ctx, id)
	if err != nil {
		return db.Order{}, nil, err
	}
	return order, items, nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID uuid.UUID, page, limit int32) ([]db.Order, int64, error) {
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	orders, err := s.repo.ListOrdersByUser(ctx, db.ListOrdersByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountOrdersByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, statusStr string) (db.Order, error) {
	var st db.OrderStatus
	switch statusStr {
	case "pending":
		st = db.OrderStatusPending
	case "confirmed":
		st = db.OrderStatusConfirmed
	case "shipped":
		st = db.OrderStatusShipped
	case "delivered":
		st = db.OrderStatusDelivered
	case "cancelled":
		st = db.OrderStatusCancelled
	default:
		return db.Order{}, errors.New("invalid status")
	}

	return s.repo.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     id,
		Status: st,
	})
}