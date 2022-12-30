package grpcclient

import (
	pb "exam/post_service/genproto/customer"
	rs "exam/post_service/genproto/review"
	"fmt"

	"exam/post_service/config"

	"google.golang.org/grpc"
)

type GrpcClientI interface {
	Review() rs.ReviewServiceClient
	Customer() pb.CustomerServiceClient
}

type GrpcClient struct {
	cfg             config.Config
	reviewService   rs.ReviewServiceClient
	customerService pb.CustomerServiceClient
}

func New(cfg config.Config) (*GrpcClient, error) {
	connCustomer, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CustomerServiceHost, cfg.CustomerServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("customer service dial host:%s, port: %d", cfg.CustomerServiceHost, cfg.CustomerServicePort)
	}
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.ReviewServiceHost, cfg.ReviewServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("review service dial host:%s, port: %d", cfg.ReviewServiceHost, cfg.ReviewServicePort)
	}
	return &GrpcClient{
		cfg:             cfg,
		reviewService:   rs.NewReviewServiceClient(connReview),
		customerService: pb.NewCustomerServiceClient(connCustomer),
	}, nil
}

func (s *GrpcClient) Review() rs.ReviewServiceClient {
	return s.reviewService
}
func (s *GrpcClient) Customer() pb.CustomerServiceClient {
	return s.customerService
}
