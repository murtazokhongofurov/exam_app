package service

import (
	l "exam/customer_service/pkg/logger"
	grpcclient "exam/customer_service/service/grpc_client"
	"exam/customer_service/storage"

	"github.com/jmoiron/sqlx"
)

type CustomerService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcclient.GrpcClientI
}

func NewUserService(db *sqlx.DB, log l.Logger, client grpcclient.GrpcClientI) *CustomerService {
	return &CustomerService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}
