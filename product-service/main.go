package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/product/v1"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/handler"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/repository"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("cannot open db: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	repo := repository.New(sqlDB)
	svc := service.New(repo)
	h := handler.New(svc)

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("product-service listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}