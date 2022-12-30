package storage

import (
	"exam/review_service/storage/postgres"
	"exam/review_service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Review() repo.ReviewStorageI
}

type storagePg struct {
	db       *sqlx.DB
	reviewResp repo.ReviewStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		reviewResp: postgres.NewReviewRepo(db),
	}
}

func (s storagePg) Review() repo.ReviewStorageI {
	return s.reviewResp
}
