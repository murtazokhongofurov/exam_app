package services

import (
	"fmt"
	"exam/api-gateway/config"
	pb "exam/api-gateway/genproto/customer"
	pbp "exam/api-gateway/genproto/post"
	rs "exam/api-gateway/genproto/review"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	CustomerService() pb.CustomerServiceClient
	PostService() pbp.PostServiceClient
	ReviewService() rs.ReviewServiceClient
}

type serviceManager struct {
	userService pb.CustomerServiceClient
	postService pbp.PostServiceClient
	reviewService rs.ReviewServiceClient
}

func (s *serviceManager) CustomerService() pb.CustomerServiceClient {
	return s.userService
}

func (s *serviceManager) PostService() pbp.PostServiceClient{
	return s.postService
}
func (s *serviceManager) ReviewService() rs.ReviewServiceClient {
	return s.reviewService
}
func NewServiceManager(conf *config.Config) (IServiceManager, error) {

	resolver.SetDefaultScheme("dns")

	connCustomer, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.CustomerServiceHost, conf.CustomerServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &serviceManager{}, err
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d",conf.PostServiceHost,conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return &serviceManager{} ,err
		}
		
	connReview,err := grpc.Dial(
		fmt.Sprintf("%s:%d",conf.ReviewServiceHost,conf.ReviewServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return &serviceManager{}, err
		}

	serviceManager := &serviceManager{
		userService: pb.NewCustomerServiceClient(connCustomer),
		postService: pbp.NewPostServiceClient(connPost),
		reviewService: rs.NewReviewServiceClient(connReview),
	}
	return serviceManager, nil
}
