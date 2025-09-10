package store_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *storeRepository) UpdateStore(context context.Context, tx *sql.Tx, store *models.Store) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:UpdateStore")
	defer span.End()

	query := `
		UPDATE stores SET
			name = $1,
			updated_at = $2
		WHERE id = $3
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	_, err := dbTx.ExecContext(context, query, store.Name, store.UpdatedAt, store.Id)

	utils.TraceError(span, err)
	return err
}
