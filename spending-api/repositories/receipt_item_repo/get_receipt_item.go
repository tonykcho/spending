package receipt_item_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

func (repo *receiptItemRepository) GetItemsByReceiptId(ctx context.Context, tx *sql.Tx, receiptId int) ([]*models.ReceiptItem, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetStoresByCategoryId")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			receipt_id,
			name,
			price,
			created_at,
			updated_at
		FROM receipt_items
		WHERE receipt_id = $1
		AND is_deleted = FALSE
		ORDER BY id
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query, receiptId)
	}

	items, err := repositories.QueryList(span, dbQuery, readReceiptItem)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	return items, nil
}

func (repo *receiptItemRepository) GetItemsByReceiptIds(ctx context.Context, tx *sql.Tx, receiptIds []int) (map[int][]*models.ReceiptItem, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetStoresByCategoryIds")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			receipt_id,
			name,
			price,
			created_at,
			updated_at
		FROM receipt_items
		WHERE receipt_id = ANY($1)
		AND is_deleted = FALSE
		ORDER BY id
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query, pq.Array(receiptIds))
	}

	items, err := repositories.QueryList(span, dbQuery, readReceiptItem)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	itemMap := make(map[int][]*models.ReceiptItem)
	for _, item := range items {
		itemMap[item.ReceiptId] = append(itemMap[item.ReceiptId], item)
	}

	return itemMap, nil
}

func readReceiptItem(rows *sql.Rows) *models.ReceiptItem {
	var item models.ReceiptItem

	err := rows.Scan(
		&item.Id,
		&item.UUId,
		&item.ReceiptId,
		&item.Name,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt)

	utils.CheckError(err)
	return &item
}
