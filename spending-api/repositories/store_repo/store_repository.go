package store_repo

import (
	"context"
	"database/sql"
	"spending/models"

	"github.com/google/uuid"
)

type StoreRepository interface {
	InsertStore(ctx context.Context, tx *sql.Tx, store *models.Store) (*models.Store, error)
	DeleteStore(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
	GetStoreById(context context.Context, tx *sql.Tx, id int) (*models.Store, error)
	GetStoreByUUId(ctx context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.Store, error)
	GetStoreByCategoryAndName(ctx context.Context, tx *sql.Tx, categoryId int, name string) (*models.Store, error)
	GetStoreList(ctx context.Context, tx *sql.Tx) ([]*models.Store, error)
}

type storeRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) StoreRepository {
	return &storeRepository{db: db}
}
