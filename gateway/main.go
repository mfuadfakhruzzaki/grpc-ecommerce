package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/gateway/middleware"
	pborder "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/order/v1"
	pbproduct "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/product/v1"
	pbuser "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	userAddr := os.Getenv("USER_SERVICE_ADDR")
	productAddr := os.Getenv("PRODUCT_SERVICE_ADDR")
	orderAddr := os.Getenv("ORDER_SERVICE_ADDR")

	if httpPort == "" {
		httpPort = "8080"
	}
	if userAddr == "" {
		userAddr = "localhost:50051"
	}
	if productAddr == "" {
		productAddr = "localhost:50052"
	}
	if orderAddr == "" {
		orderAddr = "localhost:50053"
	}

	ctx := context.Background()

	// Forward header x-user-id ke gRPC metadata
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if key == "X-User-Id" {
				return "x-user-id", true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md := metadata.MD{}
			if userID := r.Header.Get("x-user-id"); userID != "" {
				md["x-user-id"] = []string{userID}
			}
			return md
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	if err := pbuser.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userAddr, opts); err != nil {
		log.Fatalf("failed to register user service: %v", err)
	}
	if err := pbproduct.RegisterProductServiceHandlerFromEndpoint(ctx, mux, productAddr, opts); err != nil {
		log.Fatalf("failed to register product service: %v", err)
	}
	if err := pborder.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, orderAddr, opts); err != nil {
		log.Fatalf("failed to register order service: %v", err)
	}

	handler := middleware.RateLimit(middleware.JWTAuth(mux))

	log.Printf("gateway listening on :%s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, handler); err != nil {
		log.Fatalf("gateway failed: %v", err)
	}
}