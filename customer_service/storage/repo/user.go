package repo

import (
	pb "exam/customer_service/genproto/customer"
)

type CustomerStorageI interface {
	Create(*pb.CustomerRequest) (*pb.CustomerResponse, error)
	GetCustomerInfo(*pb.CustomerID) (*pb.CustomerInfo, error)
	UpdateCustomer(*pb.CustomerUpdate) (*pb.CustomerResponse, error)
	DeleteCustomer(*pb.CustomerID) (*pb.Empty, error)
	CheckFiedld(*pb.CheckFieldReq) (*pb.CheckFieldResp, error)
	GetByEmail(*pb.EmailReq) (*pb.LoginResponse, error)
	GetCustomerBySearchOrder(*pb.GetListUserRequest) (*pb.CustomerAll, error)
	GetAdmin(*pb.GetAdminReq) (*pb.GetAdminRes, error)
	GetModerator(*pb.GetModeratorReq) (*pb.GetModeratorRes, error)
}
