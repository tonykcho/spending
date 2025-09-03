package category_repo

import (
	"context"
	"database/sql"
	"spending/models"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	InsertCategory(context context.Context, category models.Category) (int, error)
	DeleteCategory(context context.Context, uuid uuid.UUID) error
	GetCategoryById(context context.Context, id int) (*models.Category, error)
	GetCategoryByUUId(context context.Context, uuid uuid.UUID) (*models.Category, error)
	GetCategoryByName(context context.Context, name string) (*models.Category, error)
	GetCategoryList(context context.Context) ([]*models.Category, error)
	GetCategoryListByIds(context context.Context, ids []int) ([]*models.Category, error)
	UpdateCategory(context context.Context, category *models.Category) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *categoryRepository {
	return &categoryRepository{db: db}
}
