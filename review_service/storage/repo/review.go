package repo

import (
	pb "exam/review_service/genproto/review"
)

type ReviewStorageI interface {
	CreateReview(*pb.ReviewRequest) (*pb.ReviewResponse, error)
	GetReviewById(*pb.ReviewId) (*pb.ReviewResponse, error)
	GetReviewPost(postId string) (*pb.Reviews, error)
	// GetReviewCustomer(ownerId string) (*pb.Reviews, error)
	UpdateReview(*pb.ReviewUp) (*pb.ReviewResponse, error)
	DeleteReview(*pb.ReviewId) (*pb.Empty, error)
	DeletePostReview(*pb.GetReviewPostRequest) (*pb.Empty, error)
	DeleteCustomerReview(*pb.CustomerDelReview) (*pb.Empty, error)
}
