package service

import (
	"context"
	pb "exam/review_service/genproto/review"
	l "exam/review_service/pkg/logger"
	"exam/review_service/storage"

	grpcclient "exam/review_service/service/grpc_client"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcclient.GrpcClientI
}

func NewReviewService(db *sqlx.DB, log l.Logger, client grpcclient.GrpcClientI) *ReviewService {
	return &ReviewService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.ReviewRequest) (*pb.ReviewResponse, error) {
	user, err := s.storage.Review().CreateReview(req)
	if err != nil {
		s.logger.Error("error insert review", l.Any("Error insert review", err))
		return &pb.ReviewResponse{}, status.Error(codes.Internal, "something went wrong, please check review info")
	}
	return user, nil
}

func (s *ReviewService) GetReviewById(ctx context.Context, req *pb.ReviewId) (*pb.ReviewResponse, error) {
	review, err := s.storage.Review().GetReviewById(req)
	if err != nil {
		s.logger.Error("error select", l.Any("Error select review", err))
		return &pb.ReviewResponse{}, status.Error(codes.Internal, "something went wrong, please check review info")
	}
	return review, nil
}

func (s *ReviewService) GetReviewPost(cix context.Context, req *pb.GetReviewPostRequest) (*pb.Reviews, error) {
	review, err := s.storage.Review().GetReviewPost(req.PostId)
	if err != nil {
		s.logger.Error("error get review post", l.Any("error getting review post", err))
		return &pb.Reviews{}, status.Error(codes.Internal, "sometging went wrong, please check review post info")
	}
	return review, nil
}

// func (s *ReviewService) GetReviewCustomer(ctx context.Context, req *pb.GetReviewCustomerRequest) (*pb.Reviews, error) {
// 	review, err := s.storage.Review().GetReviewCustomer(req.OwnerId)
// 	if err != nil {
// 		s.logger.Error("error get review customer", l.Any("error getting review customer", err))
// 		return &pb.Reviews{}, status.Error(codes.Internal, "sometging went wrong, please check review customer info")
// 	}
// 	return review, nil
// }

func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.ReviewUp) (*pb.ReviewResponse, error) {
	review, err := s.storage.Review().UpdateReview(req)
	if err != nil {
		s.logger.Error("error update", l.Any("Error update review", err))
		return &pb.ReviewResponse{}, status.Error(codes.Internal, "something went wrogn, please check review info")
	}
	return review, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.ReviewId) (*pb.Empty, error) {
	review, err := s.storage.Review().DeleteReview(req)
	if err != nil {
		s.logger.Error("error delete", l.Any("Error delete review", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrogn, please check review info")
	}
	return review, nil
}

func (s *ReviewService) DeleteCustomerReview(ctx context.Context, req *pb.CustomerDelReview) (*pb.Empty, error) {
	review, err := s.storage.Review().DeleteCustomerReview(req)
	if err != nil {
		s.logger.Error("error delete review customer", l.Any("Error delete review customer", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrogn, please check review customer info")
	}
	return review, nil
}

func (s *ReviewService) DeletePostReview(ctx context.Context, req *pb.GetReviewPostRequest) (*pb.Empty, error) {
	review, err := s.storage.Review().DeletePostReview(req)
	if err != nil {
		s.logger.Error("error delete review post", l.Any("Error delete review post", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrogn, please check review post info")
	}
	return review, nil
}
