package postgres

import (
	"database/sql"
	pb "exam/customer_service/genproto/customer"
	"fmt"
)

func (r *customerRepo) GetAdmin(req *pb.GetAdminReq) (*pb.GetAdminRes, error) {
	res := pb.GetAdminRes{}
	err := r.db.QueryRow(`SELECT 
		id,
		admin_name,
		admin_password, 
		created_at, 
		updated_at
		FROM admin 
		WHERE deleted_at 
		IS NULL AND 
		admin_name=$1`, req.Name).Scan(
		&res.Id,
		&res.Name,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		fmt.Println("Error while getting admin no rows")
		return &res, nil
	}
	if err != nil {
		return &pb.GetAdminRes{}, err
	}
	return &res, nil
}
