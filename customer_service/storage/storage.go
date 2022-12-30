package storage

import (
	"exam/customer_service/storage/postgres"
	"exam/customer_service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Customer() repo.CustomerStorageI
}

type storagePg struct {
	db       *sqlx.DB
	customerResp repo.CustomerStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		customerResp: postgres.NewCustomerRepo(db),
	}
}

func (s storagePg) Customer() repo.CustomerStorageI {
	return s.customerResp
}
