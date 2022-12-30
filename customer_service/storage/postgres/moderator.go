package postgres

import (
	"database/sql"
	pb "exam/customer_service/genproto/customer"
	"fmt"
)

func (r *customerRepo) GetModerator(req *pb.GetModeratorReq) (*pb.GetModeratorRes, error) {
	res := pb.GetModeratorRes{}
	err := r.db.QueryRow(`SELECT 
		id, 
		name, 
		password,
		created_at,
		updated_at
		FROM moderator
		WHERE deleted_at IS NULL AND name=$1`, req.Name).Scan(
		&res.Id,
		&res.Name,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt)

	if err == sql.ErrNoRows {
		fmt.Println("Error while getting admin no rows")
		return &res, nil
	}
	if err != nil {
		fmt.Println("error while getting moderator")
		return &pb.GetModeratorRes{}, err
	}
	return &res, nil
}
