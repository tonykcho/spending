package category_repo

import (
	"context"
	"database/sql"
	"spending/repositories"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *categoryRepository) DeleteCategory(context context.Context, tx *sql.Tx, uuid uuid.UUID) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:DeleteCategory")
	defer span.End()

	query := `
		UPDATE categories
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
