package service

import (
	"context"
	pb "exam/customer_service/genproto/customer"
	"exam/customer_service/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CustomerService) GetModerator(ctx context.Context, req *pb.GetModeratorReq) (*pb.GetModeratorRes, error) {
	moder, err := s.storage.Customer().GetModerator(req)
	if err != nil {
		s.logger.Error("Error while getting addmin by adminname", logger.Any("get", err))
		return &pb.GetModeratorRes{}, status.Error(codes.NotFound, "Your are not admin")
	}

	return moder, nil
}
