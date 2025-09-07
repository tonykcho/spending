package store_repo

import (
	"context"
	"database/sql"
	"spending/repositories"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *storeRepository) DeleteStore(ctx context.Context, tx *sql.Tx, uuid uuid.UUID) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:DeleteStore")
	defer span.End()

	query := `
		UPDATE stores 
		SET is_deleted = TRUE, deleted_at = NOW() 
		WHERE uuid = $1
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	_, err := dbTx.ExecContext(ctx, query, uuid)
	utils.TraceError(span, err)
	return err
}
