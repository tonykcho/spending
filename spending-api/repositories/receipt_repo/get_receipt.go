package receipt_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *receiptRepository) GetReceiptByUUId(ctx context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.Receipt, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetReceiptByUUId")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			store_name,
			total,
			date,
			created_at,
			updated_at
		FROM receipts
		WHERE uuid = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query, uuid)
	}

	receipt, err := repositories.Query(span, dbQuery, readReceipt)

	return receipt, err
}

func (repo *receiptRepository) GetReceipts(ctx context.Context, tx *sql.Tx) ([]*models.Receipt, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetReceipts")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			store_name,
			total,
			date,
			created_at,
			updated_at
		FROM receipts
		WHERE is_deleted = FALSE
		ORDER BY date DESC, created_at DESC
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query)
	}

	receipts, err := repositories.QueryList(span, dbQuery, readReceipt)

	return receipts, err
}

func readReceipt(rows *sql.Rows) *models.Receipt {
	var receipt models.Receipt

	err := rows.Scan(
		&receipt.Id,
		&receipt.UUId,
		&receipt.StoreName,
		&receipt.Total,
		&receipt.Date,
		&receipt.CreatedAt,
		&receipt.UpdatedAt)

	utils.CheckError(err)
	return &receipt
}
