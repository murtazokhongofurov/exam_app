package models

type ReviewReq struct {
	PostId      string `json:"post_id"`
	Name        string `json:"name"`
	Rating      int64  `json:"rating"`
	Description string `json:"description"`
}

type ReviewUpdate struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Rating      int64  `json:"rating"`
	Description string `json:"description"`
}
