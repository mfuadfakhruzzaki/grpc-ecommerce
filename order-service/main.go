package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/order/v1"
	pbproduct "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/product/v1"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/order-service/internal/handler"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/order-service/internal/repository"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/order-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	grpcPort := os.Getenv("GRPC_PORT")
	productAddr := os.Getenv("PRODUCT_SERVICE_ADDR")

	if grpcPort == "" {
		grpcPort = "50053"
	}
	if productAddr == "" {
		productAddr = "localhost:50052"
	}

	// Connect ke product-service
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect to product-service: %v", err)
	}
	defer productConn.Close()
	productClient := pbproduct.NewProductServiceClient(productConn)

	// Connect ke DB
	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("cannot open db: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	repo := repository.New(sqlDB)
	svc := service.New(repo, productClient)
	h := handler.New(svc)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("order-service listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}