package main

import (
	"exam/review_service/config"
	pb "exam/review_service/genproto/review"
	"exam/review_service/kafka"
	"exam/review_service/pkg/db"
	"exam/review_service/pkg/logger"
	"exam/review_service/service"
	grpcclient "exam/review_service/service/grpc_client"
	"net"
	"sync"

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

	kafkaConn, close, err := kafka.NewKafkaReader(cfg, log, connDB)
	if err != nil {
		log.Fatal("error while connection to kafka")
	}
	defer close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		kafkaConn.Reads().Start()
		wg.Done()
	}()
	
	reviewService := service.NewReviewService(connDB, log, grpcClient)
	
	lis, err := net.Listen("tcp", ":"+cfg.ReviewServicePort)
	
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	s := grpc.NewServer()
	reflection.Register(s)
	
	pb.RegisterReviewServiceServer(s, reviewService)
	
	log.Info("main: server running",
	logger.String("port", cfg.ReviewServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	wg.Wait()
}
