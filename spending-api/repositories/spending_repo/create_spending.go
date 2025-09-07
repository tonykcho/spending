package spending_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *spendingRepository) InsertSpendingRecord(ctx context.Context, tx *sql.Tx, record models.SpendingRecord) (int, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:InsertSpendingRecord")
	defer span.End()

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
		RETURNING id
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var id int
	err := dbTx.QueryRowContext(ctx, query,
		record.Amount,
		record.Remark,
		record.SpendingDate,
		record.CategoryId,
		record.CreatedAt,
		record.UpdatedAt,
	).Scan(&id)

	utils.TraceError(span, err)
	return id, err
}
