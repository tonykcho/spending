package spending_repo

import (
	"context"
	"spending/models"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *spendingRepository) InsertSpendingRecord(context context.Context, record models.SpendingRecord) (int, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:InsertSpendingRecord")
	defer span.End()

	// Create query to insert a new spending record
	query := `INSERT INTO spending_records (
				amount,
				remark,
				spending_date,
				category_id,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6)
			 RETURNING id`

	var id int
	err := repo.db.QueryRow(query,
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
