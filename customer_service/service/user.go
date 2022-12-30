package service

import (
	"context"
	pb "exam/customer_service/genproto/customer"
	pbp "exam/customer_service/genproto/post"
	"exam/customer_service/genproto/review"
	l "exam/customer_service/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CustomerService) Create(ctx context.Context, req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	user, err := s.storage.Customer().Create(req)
	if err != nil {
		s.logger.Error("error insert customers", l.Any("Error insert customers", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "something went wrong, please check customers info")
	}
	return user, nil
}

func (s *CustomerService) GetCustomerInfo(ctx context.Context, req *pb.CustomerID) (*pb.CustomerInfo, error) {
	user, err := s.storage.Customer().GetCustomerInfo(req)
	if err != nil {
		s.logger.Error("error getting customers posts", l.Any("Error getting customers posts", err))
		return &pb.CustomerInfo{}, status.Error(codes.Internal, "something went wrong, please check customers post info")
	}
	// post, err := s.client.Post().GetPostCustomer(ctx, &pbp.GetCustomerPostRequest{
	// 	OwnerId: req.Id,
	// })
	// if err != nil {
	// 	s.logger.Error("error getting  posts customer", l.Any("Error getting  posts customer", err))
	// 	return &pb.CustomerInfo{}, status.Error(codes.Internal, "something went wrong, please check  post customer info")
	// }
	// for _, p := range post.Posts {
	// 	user.Posts = append(user.Posts, &pb.Post{
	// 		Id:          p.Id,
	// 		OwnerId:     p.OwnerId,
	// 		Name:        p.Name,
	// 		Description: p.Description,
	// 	})
	// }

	// review, err := s.client.Review().GetReviewCustomer(ctx, &review.GetReviewCustomerRequest{
	// 	OwnerId: req.Id,
	// })
	// if err != nil {
	// 	s.logger.Error("error getting  review customer", l.Any("Error getting  review customer", err))
	// 	return &pb.CustomerInfo{}, status.Error(codes.Internal, "something went wrong, please check  review customer info")
	// }
	// for _, r := range review.Reviews {
	// 	user.Reviews = append(user.Reviews, &pb.ReviewResponse{
	// 		Id:          r.Id,
	// 		PostId:      r.PostId,
	// 		OwnerId:     r.OwnerId,
	// 		Name:        r.Name,
	// 		Description: r.Description,
	// 		Rating:      r.Rating,
	// 		CreatedAt:   r.CreatedAt,
	// 		UdpatedAt:   r.UdpatedAt,
	// 	})
	// }
	return user, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, req *pb.CustomerUpdate) (*pb.CustomerResponse, error) {
	user, err := s.storage.Customer().UpdateCustomer(req)
	if err != nil {
		s.logger.Error("error update", l.Any("Error update customers", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "something went wrogn, please check customer info")
	}
	return user, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, req *pb.CustomerID) (*pb.Empty, error) {
	user, err := s.storage.Customer().DeleteCustomer(req)
	if err != nil {
		s.logger.Error("error delete", l.Any("Error delete customers", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrogn, please check customer info")
	}
	_, err = s.client.Post().DeleteCustomerPost(ctx, &pbp.CustomerId{OwnerId: req.Id})
	if err != nil {
		s.logger.Error("error delete", l.Any("Error delete customers", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrogn, please check customer info")
	}
	_, err = s.client.Review().DeleteCustomerReview(ctx, &review.CustomerDelReview{OwnerId: req.Id})
	if err != nil {
		s.logger.Error("error delete", l.Any("Error delete customers", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrogn, please check customer info")
	}
	return user, nil
}

func (s *CustomerService) CheckField(ctx context.Context, req *pb.CheckFieldReq) (*pb.CheckFieldResp, error) {
	res, err := s.storage.Customer().CheckFiedld(req)
	if err != nil {
		s.logger.Error("error delete", l.Any("Error delete customers", err))
		return &pb.CheckFieldResp{}, status.Error(codes.Internal, "something went wrogn, please check customer info")

	}
	return res, nil
}

func (s *CustomerService) GetByEmail(ctx context.Context, req *pb.EmailReq) (*pb.LoginResponse, error) {
	customer, err := s.storage.Customer().GetByEmail(req)
	if err != nil {
		s.logger.Error("Error while getting customer info by email", l.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Something went wrong")
	}
	return customer, nil
}

func (s *CustomerService) GetCustomerBySearchOrder(ctx context.Context, req *pb.GetListUserRequest) (*pb.CustomerAll, error) {
	customer, err := s.storage.Customer().GetCustomerBySearchOrder(req)
	if err != nil {
		s.logger.Error("Error while getting customer info by search", l.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Something went wrong")
	}
	return customer, nil
}
