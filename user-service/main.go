package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/user/v1"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/user-service/internal/handler"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/user-service/internal/repository"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/user-service/internal/service"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
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
	pb.RegisterUserServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("user-service listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}