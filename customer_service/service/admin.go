package service

import (
	"context"
	pb "exam/customer_service/genproto/customer"
	"exam/customer_service/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CustomerService) GetAdmin(ctx context.Context, req *pb.GetAdminReq) (*pb.GetAdminRes, error) {
	res, err := s.storage.Customer().GetAdmin(req)
	if err != nil {
		s.logger.Error("Error while getting addmin by adminname", logger.Any("get", err))
		return &pb.GetAdminRes{}, status.Error(codes.NotFound, "Your are not admin")
	}
	return res, nil
}
