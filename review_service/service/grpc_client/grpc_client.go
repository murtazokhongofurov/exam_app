package grpcclient

import (
	"fmt"
	"exam/review_service/config"
	pbp "exam/review_service/genproto/post"

	"google.golang.org/grpc"
)

type GrpcClientI interface {
	Post() pbp.PostServiceClient
}

type GrpcClient struct {
	cfg         config.Config
	postService pbp.PostServiceClient
}

func New(cfg config.Config) (*GrpcClient, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("post service dial host:%s, port: %d", cfg.PostServiceHost, cfg.PostServicePort)
	}
	return &GrpcClient{
		cfg:         cfg,
		postService: pbp.NewPostServiceClient(connPost),
	}, nil
}

func (s *GrpcClient) Post() pbp.PostServiceClient {
	return s.postService
}
