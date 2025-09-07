package spending_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories/category_repo"

	"github.com/google/uuid"
)

type SpendingRepository interface {
	InsertSpendingRecord(context context.Context, tx *sql.Tx, record models.SpendingRecord) (int, error)
	GetSpendingById(context context.Context, tx *sql.Tx, id int) (*models.SpendingRecord, error)
	GetSpendingByUUId(context context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.SpendingRecord, error)
	GetSpendingList(context context.Context, tx *sql.Tx) ([]*models.SpendingRecord, error)
	LoadSpendingCategory(context context.Context, tx *sql.Tx, record *models.SpendingRecord) error
	LoadSpendingListCategory(context context.Context, tx *sql.Tx, records []*models.SpendingRecord) error
	DeleteSpending(context context.Context, tx *sql.Tx, uuid uuid.UUID) error
}

type spendingRepository struct {
	db            *sql.DB
	category_repo category_repo.CategoryRepository
}

func NewSpendingRepository(db *sql.DB) *spendingRepository {
	return &spendingRepository{db: db, category_repo: category_repo.NewCategoryRepository(db)}
}
