package receipt_repo

import (
	"context"
	"database/sql"
	"spending/repositories"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *receiptRepository) DeleteReceipt(context context.Context, tx *sql.Tx, uuid uuid.UUID) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:DeleteReceipt")
	defer span.End()

	query := `
		UPDATE receipts
		SET is_deleted = TRUE,
			deleted_at = NOW()
		WHERE uuid = $1
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	_, err := dbTx.ExecContext(context, query, uuid)
	utils.TraceError(span, err)
	return err
}
