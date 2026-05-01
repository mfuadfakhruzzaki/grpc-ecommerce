package handler

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	svc "github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/service"
	pb "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	svc *svc.ProductService
}

func New(svc *svc.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductReq) (*pb.CreateProductRes, error) {
	p, err := h.svc.Create(ctx, req.Name, req.Description, req.Price, req.CategoryId, req.StockQty)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create failed: %v", err)
	}
	price, _ := strconv.ParseFloat(p.Price, 64)
	return &pb.CreateProductRes{Product: &pb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description.String,
		Price:       price,
		CategoryId:  p.CategoryID.UUID.String(),
		StockQty:    p.StockQty.Int32,
		IsActive:    p.IsActive.Bool,
	}}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.GetProductRes, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	p, err := h.svc.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}
	price, _ := strconv.ParseFloat(p.Price, 64)
	return &pb.GetProductRes{Product: &pb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description.String,
		Price:       price,
		CategoryId:  p.CategoryID.UUID.String(),
		StockQty:    p.StockQty.Int32,
		IsActive:    p.IsActive.Bool,
	}}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsRes, error) {
	products, total, err := h.svc.List(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	var pbProducts []*pb.Product
	for _, p := range products {
		price, _ := strconv.ParseFloat(p.Price, 64)
		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description.String,
			Price:       price,
			CategoryId:  p.CategoryID.UUID.String(),
			StockQty:    p.StockQty.Int32,
			IsActive:    p.IsActive.Bool,
		})
	}
	return &pb.ListProductsRes{Products: pbProducts, Total: int32(total)}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductReq) (*pb.UpdateProductRes, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	p, err := h.svc.Update(ctx, id, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}
	price, _ := strconv.ParseFloat(p.Price, 64)
	return &pb.UpdateProductRes{Product: &pb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description.String,
		Price:       price,
		CategoryId:  p.CategoryID.UUID.String(),
		StockQty:    p.StockQty.Int32,
		IsActive:    p.IsActive.Bool,
	}}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.DeleteProductRes, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	return &pb.DeleteProductRes{Success: true}, nil
}

func (h *ProductHandler) CheckStock(ctx context.Context, req *pb.CheckStockReq) (*pb.CheckStockRes, error) {
	id, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}
	available, current, err := h.svc.CheckStock(ctx, id, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.CheckStockRes{Available: available, CurrentStock: current}, nil
}

func (h *ProductHandler) DeductStock(ctx context.Context, req *pb.DeductStockReq) (*pb.DeductStockRes, error) {
	id, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}
	success, remaining, err := h.svc.DeductStock(ctx, id, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}
	return &pb.DeductStockRes{Success: success, RemainingStock: remaining}, nil
}