package postgres

import (
	"database/sql"
	pb "exam/review_service/genproto/review"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type reviewRepo struct {
	db *sqlx.DB
}

func NewReviewRepo(db *sqlx.DB) *reviewRepo {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) CreateReview(review *pb.ReviewRequest) (*pb.ReviewResponse, error) {
	reviewResp := pb.ReviewResponse{}
	err := r.db.QueryRow(`
	INSERT INTO review (
		id,
		post_id,
		owner_id,
		name,
		rating,
		description
		) 
	VALUES ($1,$2,$3,$4,$5,$6) 
		RETURNING 
		id, post_id, owner_id, name, rating, description, created_at, updated_at`,
		review.Id, review.PostId, review.OwnerId, review.Name, review.Rating, review.Description,
	).Scan(
		&reviewResp.Id,
		&reviewResp.PostId,
		&reviewResp.OwnerId,
		&reviewResp.Name,
		&reviewResp.Rating,
		&reviewResp.Description,
		&reviewResp.CreatedAt,
		&reviewResp.UdpatedAt,
	)
	if err != nil {
		return &pb.ReviewResponse{}, err
	}
	return &reviewResp, nil
}

func (r *reviewRepo) GetReviewById(req *pb.ReviewId) (*pb.ReviewResponse, error) {
	response := pb.ReviewResponse{}
	err := r.db.QueryRow(`SELECT 
			id, 
			post_id, 
			owner_id, 
			name, rating, 
			description,
			created_at,
			updated_at
			FROM review 
			WHERE id=$1 AND deleted_at IS NULL`, req.Id).Scan(
		&response.Id,
		&response.PostId,
		&response.OwnerId,
		&response.Name,
		&response.Rating,
		&response.Description,
		&response.CreatedAt,
		&response.UdpatedAt,
	)
	if err == sql.ErrNoRows {
		return &pb.ReviewResponse{}, err
	}
	if err != nil {
		fmt.Println("error whilr getting review")
		return &pb.ReviewResponse{}, err
	}
	return &response, nil
}

func (r *reviewRepo) GetReviewPost(postId string) (*pb.Reviews, error) {
	rows, err := r.db.Query(`SELECT 
		id, 
		post_id, 
		owner_id, 
		name, 
		rating, 
		description, 
		created_at, 
		updated_at 
		FROM  review WHERE post_id=$1 AND deleted_at IS NULL`, postId)
	if err != nil {
		fmt.Println("error while getting review with postId")
		return &pb.Reviews{}, err
	}
	defer rows.Close()
	reviews := pb.Reviews{}
	for rows.Next() {
		review := pb.ReviewResponse{}
		err = rows.Scan(&review.Id, &review.PostId,
			&review.OwnerId, &review.Name, &review.Rating, &review.Description,
			&review.CreatedAt, &review.UdpatedAt)
		if err != nil {
			fmt.Println("error while scanning review with postId")
			return &pb.Reviews{}, err
		}
		reviews.Reviews = append(reviews.Reviews, &review)
	}
	return &reviews, nil
}

func (r *reviewRepo) UpdateReview(req *pb.ReviewUp) (*pb.ReviewResponse, error) {
	review := pb.ReviewResponse{}
	_, err := r.db.Exec(`UPDATE review SET 
			name=$1,
			description=$2,
			rating=$3,
			updated_at=NOW()
			WHERE id = $4 AND deleted_at IS NULL`,
		req.Name, req.Description, req.Rating, req.Id)
	if err != nil {
		return &pb.ReviewResponse{}, err
	}
	err = r.db.QueryRow(`SELECT 
	id, 
	post_id,
	owner_id,
	name, 
	description, 
	rating,
	created_at,
	updated_at
	FROM review 
	WHERE id=$1 AND deleted_at IS NULL`, req.Id).Scan(
		&review.Id,
		&review.PostId,
		&review.OwnerId,
		&review.Name,
		&review.Description,
		&review.Rating,
		&review.CreatedAt,
		&review.UdpatedAt,
	)
	if err != nil {
		fmt.Println("error while getting review in update")
		return &pb.ReviewResponse{}, err
	}
	return &review, nil
}

func (r *reviewRepo) DeleteReview(req *pb.ReviewId) (*pb.Empty, error) {
	_, err := r.db.Exec(`UPDATE review SET deleted_at=NOW() WHERE deleted_at IS NULL  AND id=$1`, req.Id)
	if err != nil {
		return &pb.Empty{}, err
	}
	return &pb.Empty{}, nil
}

func (r *reviewRepo) DeleteCustomerReview(review *pb.CustomerDelReview) (*pb.Empty, error) {
	_, err := r.db.Exec(`UPDATE review SET deleted_at=NOW() WHERE deleted_at IS NULL  AND owner_id=$1`, review.OwnerId)
	if err != nil {
		fmt.Println("error while delete review customer")
		return &pb.Empty{}, err
	}
	return &pb.Empty{}, nil
}

func (r *reviewRepo) DeletePostReview(review *pb.GetReviewPostRequest) (*pb.Empty, error) {
	_, err := r.db.Exec(`UPDATE review SET deleted_at=NOW() WHERE deleted_at IS NULL AND post_id=$1`, review.PostId)
	if err != nil {
		fmt.Println("error while delete review post")
		return &pb.Empty{}, err
	}
	return &pb.Empty{}, nil
}
