package service

import (
	"exam/post_service/kafka"
	l "exam/post_service/pkg/logger"
	grpcclient "exam/post_service/service/grpc_client"
	"exam/post_service/storage"

	"github.com/jmoiron/sqlx"
)

type PostService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcclient.GrpcClientI
	Kafka   kafka.KafkaI
}

func NewPostService(db *sqlx.DB, log l.Logger, client grpcclient.GrpcClientI, kafka kafka.KafkaI) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
		Kafka:   kafka,
	}
}
