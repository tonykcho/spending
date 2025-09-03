package spending_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories/category_repo"

	"github.com/google/uuid"
)

type SpendingRepository interface {
	InsertSpendingRecord(context context.Context, record models.SpendingRecord) (int, error)
	GetSpendingById(context context.Context, id int) (*models.SpendingRecord, error)
	GetSpendingByUUId(context context.Context, uuid uuid.UUID) (*models.SpendingRecord, error)
	GetSpendingList(context context.Context) ([]*models.SpendingRecord, error)
	LoadSpendingCategory(context context.Context, record *models.SpendingRecord) error
	LoadSpendingListCategory(context context.Context, records []*models.SpendingRecord) error
}

type spendingRepository struct {
	db            *sql.DB
	category_repo category_repo.CategoryRepository
}

func NewSpendingRepository(db *sql.DB) *spendingRepository {
	return &spendingRepository{db: db, category_repo: category_repo.NewCategoryRepository(db)}
}
