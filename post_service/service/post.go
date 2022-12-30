package service

import (
	"context"

	pb "exam/post_service/genproto/post"
	rs "exam/post_service/genproto/review"
	l "exam/post_service/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PostService) CreatePost(ctx context.Context, req *pb.PostReq) (*pb.PostResp, error) {
	res, err := s.storage.Post().CreatePost(req)
	if err != nil {
		s.logger.Error("error produce review", l.Any("Error produce  review", err))
		return &pb.PostResp{}, status.Error(codes.Internal, "something went wrong, please check produce review info")
	}

	for _, r := range req.Reviews {
		review := &pb.Review{
			Id:          r.Id,
			PostId:      r.PostId,
			OwnerId:     r.OwnerId,
			Name:        r.Name,
			Rating:      r.Rating,
			Description: r.Description,
		}
		res.Reviews = append(res.Reviews, review)
	}
	s.Kafka.Produce().SendMessage(req)
	// err := s.Kafka.Produce().SendMessage(req)
	if err != nil {
		s.logger.Error("error produce review", l.Any("Error produce  review", err))
		return &pb.PostResp{}, status.Error(codes.Internal, "something went wrong, please check produce review info")
	}

	return res, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.PostUp) (*pb.PostResp, error) {
	post, err := s.storage.Post().UpdatePost(req)
	if err != nil {
		s.logger.Error("error update post", l.Any("error updating post", err))
		return &pb.PostResp{}, status.Error(codes.Internal, "sometging went wrong, please check post info")
	}
	return post, nil
}

func (s *PostService) GetPostReview(ctx context.Context, req *pb.ID) (*pb.PostInfo, error) {
	post, err := s.storage.Post().GetPostReview(req)
	if err != nil {
		s.logger.Error("error select", l.Any("Error select post review", err))
		return &pb.PostInfo{}, status.Error(codes.Internal, "something went wrong, please check post info")
	}
	review, err := s.client.Review().GetReviewPost(ctx, &rs.GetReviewPostRequest{
		PostId: req.PostID,
	})
	if err != nil {
		s.logger.Error("error select", l.Any("Error select post review", err))
		return &pb.PostInfo{}, status.Error(codes.Internal, "something went wrong, please check post with review info")
	}
	for _, r := range review.Reviews {
		post.Reviews = append(post.Reviews, &pb.Review{
			Id:          r.Id,
			PostId:      r.PostId,
			OwnerId:     r.OwnerId,
			Name:        r.Name,
			Rating:      r.Rating,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
			UdpatedAt:   r.UdpatedAt,
		})
	}
	return post, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.ID) (*pb.Empty, error) {
	post, err := s.storage.Post().DeletePost(req)
	if err != nil {
		s.logger.Error("error update post", l.Any("error updating post", err))
		return &pb.Empty{}, status.Error(codes.Internal, "sometging went wrong, please check post info")
	}
	_, err = s.client.Review().DeletePostReview(ctx, &rs.GetReviewPostRequest{PostId: req.PostID})
	if err != nil {
		return &pb.Empty{}, err
	}
	return post, nil
}

func (s *PostService) DeleteCustomerPost(ctx context.Context, req *pb.CustomerId) (*pb.Empty, error) {
	post, err := s.storage.Post().DeleteCustomerPost(req)
	if err != nil {
		s.logger.Error("error delete post", l.Any("error delete post customer", err))
		return &pb.Empty{}, status.Error(codes.Internal, "sometging went wrong, please check post info")
	}
	return post, nil
}

func (s *PostService) ListPost(ctx context.Context, req *pb.ListPostRequest) (*pb.ListPostResponse, error) {
	post, err := s.storage.Post().ListPost(req.Value, req.Limit, req.Page)
	if err != nil {
		s.logger.Error("error get  post list", l.Any("error get post list", err))
		return &pb.ListPostResponse{}, status.Error(codes.Internal, "sometging went wrong, please check post info")
	}
	return post, nil
}
