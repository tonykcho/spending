package store_repo

import "database/sql"

type StoreRepository interface {
}

type storeRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) StoreRepository {
	return &storeRepository{db: db}
}
