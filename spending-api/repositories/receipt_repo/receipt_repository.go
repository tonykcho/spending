package receipt_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories/receipt_item_repo"

	"github.com/google/uuid"
)

type ReceiptRepository interface {
	GetReceiptByUUId(ctx context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.Receipt, error)
}

type receiptRepository struct {
	db                *sql.DB
	receipt_item_repo receipt_item_repo.ReceiptItemRepository
}

func NewReceiptRepository(db *sql.DB, receiptItemRepo receipt_item_repo.ReceiptItemRepository) *receiptRepository {
	return &receiptRepository{db: db, receipt_item_repo: receiptItemRepo}
}
