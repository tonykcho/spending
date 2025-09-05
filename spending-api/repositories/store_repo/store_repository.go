package store_repo

import (
	"context"
	"database/sql"
	"spending/models"

	"github.com/google/uuid"
)

type StoreRepository interface {
	InsertStore(ctx context.Context, store *models.Store) (int, error)
	DeleteStore(ctx context.Context, id uuid.UUID) error
	GetStoreById(context context.Context, id int) (*models.Store, error)
	GetStoreByUUId(ctx context.Context, uuid uuid.UUID) (*models.Store, error)
	GetStoreByCategoryAndName(ctx context.Context, categoryId int, name string) (*models.Store, error)
	GetStoreList(ctx context.Context) ([]*models.Store, error)
}

type storeRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) StoreRepository {
	return &storeRepository{db: db}
}
