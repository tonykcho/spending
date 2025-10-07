package receipt_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *receiptRepository) InsertReceipt(context context.Context, tx *sql.Tx, receipt *models.Receipt) (*models.Receipt, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:InsertReceipt")
	defer span.End()

	if receipt == nil {
		return nil, fmt.Errorf("receipt cannot be nil")
	}

	query := `
	INSERT INTO receipts (
		store_name,
		date,
		total,
		created_at,
		updated_at
	) Values ($1, $2, $3, $4, $5)
		RETURNING
			id,
			uuid,
			store_name,
			total,
			date,
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
			receipt.StoreName,
			receipt.Date,
			receipt.Total,
			receipt.CreatedAt,
			receipt.UpdatedAt)
	}

	newReceipt, err := repositories.Query(span, dbQuery, readReceipt)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	return newReceipt, nil
}
