package handler

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/user/v1"
	svc "github.com/mfuadfakhruzzaki/grpc-ecommerce/user-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	svc *svc.UserService
}

func New(svc *svc.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func userIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.UUID{}, status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("x-user-id")
	if len(vals) == 0 {
		return uuid.UUID{}, status.Error(codes.Unauthenticated, "missing user id")
	}
	id, err := uuid.Parse(vals[0])
	if err != nil {
		return uuid.UUID{}, status.Error(codes.Unauthenticated, "invalid user id")
	}
	return id, nil
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserReq) (*pb.RegisterUserRes, error) {
	user, err := h.svc.Register(ctx, req.Email, req.Password, req.FullName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}
	return &pb.RegisterUserRes{
		Id:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName.String,
	}, nil
}

func (h *UserHandler) LoginUser(ctx context.Context, req *pb.LoginUserReq) (*pb.LoginUserRes, error) {
	token, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}
	return &pb.LoginUserRes{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *pb.GetProfileReq) (*pb.GetProfileRes, error) {
	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	user, err := h.svc.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return &pb.GetProfileRes{
		Id:        user.ID.String(),
		Email:     user.Email,
		FullName:  user.FullName.String,
		AvatarUrl: user.AvatarUrl.String,
	}, nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileReq) (*pb.UpdateProfileRes, error) {
	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	user, err := h.svc.UpdateProfile(ctx, userID, req.FullName, req.AvatarUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}
	return &pb.UpdateProfileRes{
		Id:        user.ID.String(),
		FullName:  user.FullName.String,
		AvatarUrl: user.AvatarUrl.String,
	}, nil
}