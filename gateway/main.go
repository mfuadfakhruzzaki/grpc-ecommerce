package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/gateway/middleware"
	pborder "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/order/v1"
	pbproduct "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/product/v1"
	pbuser "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register user service
	if err := pbuser.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userAddr, opts); err != nil {
		log.Fatalf("failed to register user service: %v", err)
	}

	// Register product service
	if err := pbproduct.RegisterProductServiceHandlerFromEndpoint(ctx, mux, productAddr, opts); err != nil {
		log.Fatalf("failed to register product service: %v", err)
	}

	// Register order service
	if err := pborder.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, orderAddr, opts); err != nil {
		log.Fatalf("failed to register order service: %v", err)
	}

	// Chain middleware
	handler := middleware.RateLimit(middleware.JWTAuth(mux))

	log.Printf("gateway listening on :%s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, handler); err != nil {
		log.Fatalf("gateway failed: %v", err)
	}
}