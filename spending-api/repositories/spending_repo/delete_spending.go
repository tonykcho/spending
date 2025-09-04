package spending_repo

import (
	"context"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *spendingRepository) DeleteSpending(context context.Context, uuid uuid.UUID) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:DeleteSpending")
	defer span.End()

	query := `
		UPDATE spending_records
		SET is_deleted = TRUE, deleted_at = NOW()
		WHERE uuid = $1
	`

	_, err := repo.db.Exec(query, uuid)
	utils.TraceError(span, err)
	return err
}
