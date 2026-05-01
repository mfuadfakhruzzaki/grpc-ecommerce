package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/mfuadfakhruzzaki/grpc-ecommerce/proto/product/v1"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/handler"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/repository"
	"github.com/mfuadfakhruzzaki/grpc-ecommerce/product-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	repo := repository.New(sqlDB)
	svc := service.New(repo)
	h := handler.New(svc)

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(500),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     30 * time.Second,
			MaxConnectionAge:      2 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  10 * time.Second,
			Timeout:               3 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)
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