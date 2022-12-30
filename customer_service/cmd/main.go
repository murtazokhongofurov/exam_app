package main

import (
	"exam/customer_service/config"
	pb "exam/customer_service/genproto/customer"
	"exam/customer_service/pkg/db"
	"exam/customer_service/pkg/logger"
	"exam/customer_service/service"
	grpcclient "exam/customer_service/service/grpc_client"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "golang")
	defer logger.Cleanup(log)
	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase),
	)
	connDB, err := db.ConnectToDB(cfg)

	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}
	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		log.Fatal("grpc connection to client error", logger.Error(err))
	}

	userService := service.NewUserService(connDB, log, grpcClient)

	lis, err := net.Listen("tcp", ":"+cfg.CustomerServicePort)

	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterCustomerServiceServer(s, userService)

	log.Info("main: server running",
		logger.String("port", cfg.CustomerServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
