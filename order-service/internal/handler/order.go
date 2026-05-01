package handler

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	pb "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/order/v1"
	svc "github.com/mfuadfakhruzzaki/grpc-ecommerce/order-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	svc *svc.OrderService
}

func New(svc *svc.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderRes, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	var items []svc.OrderItemInput
	for _, item := range req.Items {
		items = append(items, svc.OrderItemInput{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	order, orderItems, err := h.svc.CreateOrder(ctx, userID, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create order failed: %v", err)
	}

	totalAmount, _ := strconv.ParseFloat(order.TotalAmount, 64)
	pbOrder := &pb.Order{
		Id:          order.ID.String(),
		UserId:      order.UserID.String(),
		Status:      string(order.Status),
		TotalAmount: totalAmount,
	}
	for _, item := range orderItems {
		price, _ := strconv.ParseFloat(item.PriceAtOrder, 64)
		pbOrder.Items = append(pbOrder.Items, &pb.OrderItemDetail{
			ProductId:    item.ProductID.String(),
			Quantity:     item.Quantity,
			PriceAtOrder: price,
		})
	}

	return &pb.CreateOrderRes{Order: pbOrder}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderReq) (*pb.GetOrderRes, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	order, items, err := h.svc.GetOrder(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	totalAmount, _ := strconv.ParseFloat(order.TotalAmount, 64)
	pbOrder := &pb.Order{
		Id:          order.ID.String(),
		UserId:      order.UserID.String(),
		Status:      string(order.Status),
		TotalAmount: totalAmount,
	}
	for _, item := range items {
		price, _ := strconv.ParseFloat(item.PriceAtOrder, 64)
		pbOrder.Items = append(pbOrder.Items, &pb.OrderItemDetail{
			ProductId:    item.ProductID.String(),
			Quantity:     item.Quantity,
			PriceAtOrder: price,
		})
	}

	return &pb.GetOrderRes{Order: pbOrder}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersReq) (*pb.ListOrdersRes, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	orders, total, err := h.svc.ListOrders(ctx, userID, req.Page, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list orders failed: %v", err)
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		totalAmount, _ := strconv.ParseFloat(order.TotalAmount, 64)
		pbOrders = append(pbOrders, &pb.Order{
			Id:          order.ID.String(),
			UserId:      order.UserID.String(),
			Status:      string(order.Status),
			TotalAmount: totalAmount,
		})
	}

	return &pb.ListOrdersRes{Orders: pbOrders, Total: int32(total)}, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateStatusReq) (*pb.UpdateStatusRes, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	order, err := h.svc.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update status failed: %v", err)
	}

	totalAmount, _ := strconv.ParseFloat(order.TotalAmount, 64)
	return &pb.UpdateStatusRes{Order: &pb.Order{
		Id:          order.ID.String(),
		UserId:      order.UserID.String(),
		Status:      string(order.Status),
		TotalAmount: totalAmount,
	}}, nil
}