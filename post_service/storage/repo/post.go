package repo

import (
	pb "exam/post_service/genproto/post"
)

type PostStorageI interface {
	CreatePost(*pb.PostReq) (*pb.PostResp, error)
	GetPostReview(*pb.ID) (*pb.PostInfo, error)
	// GetPostCustomer(OwnerId string) (*pb.Posts, error)
	UpdatePost(*pb.PostUp) (*pb.PostResp, error)
	DeletePost(*pb.ID) (*pb.Empty, error)
	DeleteCustomerPost(*pb.CustomerId) (*pb.Empty, error)
	ListPost(value string, limit, page int64) (*pb.ListPostResponse, error)
}
