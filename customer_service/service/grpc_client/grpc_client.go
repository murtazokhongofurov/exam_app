package grpcclient

import (
	"exam/customer_service/config"
	pbp "exam/customer_service/genproto/post"
	rs "exam/customer_service/genproto/review"
	"fmt"

	"google.golang.org/grpc"
)

type GrpcClientI interface {
	Post() pbp.PostServiceClient
	Review() rs.ReviewServiceClient
}

type GrpcClient struct {
	cfg           config.Config
	postService   pbp.PostServiceClient
	reviewService rs.ReviewServiceClient
}

func New(cfg config.Config) (*GrpcClient, error) {
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.ReviewServiceHost, cfg.ReviewServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("review service dial host:%s, port: %d", cfg.ReviewServiceHost, cfg.ReviewServicePort)
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("post service dial host:%s, port: %d", cfg.PostServiceHost, cfg.PostServicePort)
	}
	return &GrpcClient{
		cfg:           cfg,
		postService:   pbp.NewPostServiceClient(connPost),
		reviewService: rs.NewReviewServiceClient(connReview),
	}, nil
}

func (s *GrpcClient) Review() rs.ReviewServiceClient {
	return s.reviewService
}
func (s *GrpcClient) Post() pbp.PostServiceClient {
	return s.postService
}
