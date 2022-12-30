package main

import (
	"exam/post_service/config"
	pb "exam/post_service/genproto/post"
	"exam/post_service/kafka"
	"exam/post_service/pkg/db"
	"exam/post_service/pkg/logger"
	"exam/post_service/service"
	grpcclient "exam/post_service/service/grpc_client"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "golang")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PosgresHost),
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

	kafka, close, err := kafka.NewKafka(cfg)
	if err != nil {
		log.Fatal("Error while connecting to kafka", logger.Error(err))
	}
	defer close()
	postService := service.NewPostService(connDB, log, grpcClient, kafka)
	lis, err := net.Listen("tcp", ":"+cfg.PostServicePort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterPostServiceServer(s, postService)
	log.Info("main:server running",
		logger.String("port", cfg.PostServicePort),
	)
	if err := s.Serve(lis); err != nil {
		log.Fatal("error while listening2: %v", logger.Error(err))
	}
}
