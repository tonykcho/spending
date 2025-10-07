package receipt_item_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *receiptItemRepository) InsertReceiptItem(context context.Context, tx *sql.Tx, receiptItem *models.ReceiptItem) (*models.ReceiptItem, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:InsertReceiptItem")
	defer span.End()

	if receiptItem == nil {
		return nil, fmt.Errorf("receipt item cannot be nil")
	}

	query := `
	INSERT INTO receipt_items (
		receipt_id,
		name,
		price,
		created_at,
		updated_at
	) Values ($1, $2, $3, $4, $5)
		RETURNING
			id,
			uuid,
			receipt_id,
			name,
			price,
			created_at,
			updated_at
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.QueryContext(context,
			query,
			receiptItem.ReceiptId,
			receiptItem.Name,
			receiptItem.Price,
			receiptItem.CreatedAt,
			receiptItem.UpdatedAt)
	}

	newReceiptItem, err := repositories.Query(span, dbQuery, readReceiptItem)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	return newReceiptItem, nil
}
