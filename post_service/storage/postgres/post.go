package postgres

import (
	"database/sql"
	pb "exam/post_service/genproto/post"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(post *pb.PostReq) (*pb.PostResp, error) {
	tr, _ := r.db.Begin()
	defer tr.Rollback()
	responsePost := pb.PostResp{}
	err := tr.QueryRow(`INSERT INTO posts(
		id,
		owner_id,
		name,
		description) 
		VALUES($1, $2, $3, $4) 
		RETURNING 
		id, 
		owner_id, 
		name,
		description,
		created_at,
		updated_at`,
		post.Id,
		post.OwnerId,
		post.Name,
		post.Description).Scan(
		&responsePost.Id,
		&responsePost.OwnerId,
		&responsePost.Name,
		&responsePost.Description,
		&responsePost.CreatedAt,
		&responsePost.UpdatedAt)
	if err != nil {
		tr.Rollback()
		fmt.Println("error while inserting posts")
		return &pb.PostResp{}, err
	}
	var medias []*pb.Media
	for _, media := range post.Medias {
		mediaResp := pb.Media{}
		media.PostId = responsePost.Id
		err = tr.QueryRow(`INSERT INTO medias(
			id, 
			post_id, 
			name, 
			link, 
			type)
			values($1,$2,$3,$4,$5) 
			RETURNING 
			id, 
			post_id, 
			name, 
			link, 
			type`,
			media.Id,
			responsePost.Id,
			media.Name,
			media.Link,
			media.Type).Scan(
			&mediaResp.Id,
			&mediaResp.PostId,
			&mediaResp.Name,
			&mediaResp.Link,
			&mediaResp.Type,
		)
		if err != nil {
			tr.Rollback()
			fmt.Println("error while inserting media")
			return &pb.PostResp{}, err
		}
		medias = append(medias, &mediaResp)
	}
	responsePost.Medias = medias

	if err = tr.Commit(); err != nil {
		fmt.Println("error tr.Commit()", err)
	}
	return &responsePost, nil
}

func (r *postRepo) GetPostReview(post *pb.ID) (*pb.PostInfo, error) {
	response := pb.PostInfo{}
	err := r.db.QueryRow(`SELECT  
		id, 
		owner_id, 
		name, 
		description,
		created_at, 
		updated_at 
		FROM  posts WHERE 
		id=$1 AND deleted_at IS NULL`, post.PostID).Scan(
		&response.Id,
		&response.OwnerId,
		&response.Name,
		&response.Description,
		&response.CreatedAt,
		&response.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return &pb.PostInfo{}, err
	}
	if err != nil {
		fmt.Println("error while getting post with review")
		return &pb.PostInfo{}, err
	}
	rows, err := r.db.Query(`SELECT 
	id, 
	post_id, 
	name, 
	link, 
	type 
	FROM medias 
	WHERE post_id=$1`, post.PostID)
	if err == sql.ErrNoRows {
		return &pb.PostInfo{}, err
	}
	if err != nil {
		fmt.Println("error while getting media")
		return &pb.PostInfo{}, err
	}
	defer rows.Close()
	for rows.Next() {
		media := pb.Media{}
		err = rows.Scan(
			&media.Id,
			&media.PostId,
			&media.Name,
			&media.Link,
			&media.Type)
		if err != nil {
			fmt.Println("error while scanning media")
			return &pb.PostInfo{}, err
		}
		response.Medias = append(response.Medias, &media)
	}
	return &response, nil
}

func (r *postRepo) UpdatePost(post *pb.PostUp) (*pb.PostResp, error) {
	response := pb.PostResp{}
	_, err := r.db.Exec(`UPDATE posts SET
		name=$1,
		description=$2,
		updated_at=NOW()
		WHERE id=$3 AND deleted_at IS NULL`, post.Name, post.Description, post.Id)
	if err != nil {
		fmt.Println("error while update post")
		return &pb.PostResp{}, err
	}
	for _, media := range post.Medias {
		_, err = r.db.Exec(`UPDATE medias SET name=$1, link=$2, type=$3 WHERE id=$4`, media.Name, media.Link, media.Type, media.Id)
		if err != nil {
			fmt.Println("error while update media")
			return &pb.PostResp{}, err
		}
	}
	return &response, nil
}

func (r *postRepo) DeletePost(post *pb.ID) (*pb.Empty, error) {
	response := pb.Empty{}
	err := r.db.QueryRow(`UPDATE posts SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, post.PostID).Err()
	if err != nil {
		fmt.Println("error delete posts")
		return &pb.Empty{}, err
	}
	return &response, nil
}

func (r *postRepo) DeleteCustomerPost(post *pb.CustomerId) (*pb.Empty, error) {
	response := pb.Empty{}
	err := r.db.QueryRow(`UPDATE posts SET deleted_at=NOW() WHERE owner_id=$1 AND deleted_at IS NULL`, post.OwnerId).Err()
	if err != nil {
		fmt.Println("error delete posts customers")
		return &pb.Empty{}, err
	}
	return &response, nil
}

func (r *postRepo) ListPost(value string, limit, page int64) (*pb.ListPostResponse, error) {
	posts := pb.ListPostResponse{}
	query := fmt.Sprintf("SELECT name, description FROM posts WHERE name LIKE '%s%%' OR description LIKE '%s%%' ORDER BY name  LIMIT %d OFFSET %d", value, value, limit, (page-1)*limit)
	rowsPost, err := r.db.Query(query)
	for rowsPost.Next() {
		post := pb.PostSearchResp{}
		err = rowsPost.Scan(
			&post.Name,
			&post.Description,
		)

		if err != nil {
			fmt.Println("error while scan post list")
			return &pb.ListPostResponse{}, err
		}
		posts.Posts = append(posts.Posts, &post)
	}
	return &posts, err
}
