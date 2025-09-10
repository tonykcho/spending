package spending_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *spendingRepository) InsertSpendingRecord(ctx context.Context, tx *sql.Tx, record *models.SpendingRecord) (*models.SpendingRecord, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:InsertSpendingRecord")
	defer span.End()

	if record == nil {
		return nil, fmt.Errorf("record is nil")
	}

	// Create query to insert a new spending record
	query := `
	INSERT INTO spending_records (
		amount,
		remark,
		spending_date,
		category_id,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			uuid,
			amount,
			remark,
			spending_date,
			category_id,
			created_at,
			updated_at
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query,
			record.Amount,
			record.Remark,
			record.SpendingDate,
			record.CategoryId,
			record.CreatedAt,
			record.UpdatedAt,
		)
	}

	newRecord, err := repositories.Query(span, dbQuery, readSpendingRecord)

	utils.TraceError(span, err)
	return newRecord, err
}
