package category_repo

import (
	"context"
	"database/sql"
	"spending/models"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	InsertCategory(context context.Context, tx *sql.Tx, category models.Category) (int, error)
	DeleteCategory(context context.Context, tx *sql.Tx, uuid uuid.UUID) error
	GetCategoryById(context context.Context, tx *sql.Tx, id int) (*models.Category, error)
	GetCategoryByUUId(context context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.Category, error)
	GetCategoryByName(context context.Context, tx *sql.Tx, name string) (*models.Category, error)
	GetCategoryList(context context.Context, tx *sql.Tx) ([]*models.Category, error)
	GetCategoryListByIds(context context.Context, tx *sql.Tx, ids []int) ([]*models.Category, error)
	UpdateCategory(context context.Context, tx *sql.Tx, category *models.Category) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *categoryRepository {
	return &categoryRepository{db: db}
}
