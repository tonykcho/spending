package receipt_item_repo

import (
	"context"
	"database/sql"
	"spending/models"
)

type ReceiptItemRepository interface {
	// Define methods for the ReceiptItemRepository here
	GetItemsByReceiptId(ctx context.Context, tx *sql.Tx, receiptId int) ([]*models.ReceiptItem, error)
	GetItemsByReceiptIds(ctx context.Context, tx *sql.Tx, receiptIds []int) (map[int][]*models.ReceiptItem, error)
}

type receiptItemRepository struct {
	db *sql.DB
}

func NewReceiptItemRepository(db *sql.DB) *receiptItemRepository {
	return &receiptItemRepository{db: db}
}
