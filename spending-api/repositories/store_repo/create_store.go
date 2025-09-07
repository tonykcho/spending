package store_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/models"
	"spending/repositories"

	"go.opentelemetry.io/otel"
)

func (repo *storeRepository) InsertStore(ctx context.Context, tx *sql.Tx, store *models.Store) (*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:InsertStore")
	defer span.End()

	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}

	query := `INSERT INTO stores (
				name,
				category_id,
				created_at,
				updated_at
			) Values ($1, $2, $3, $4)
			 RETURNING
			 	id,
				uuid,
				name,
				category_id,
				created_at,
				updated_at`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query,
			store.Name,
			store.CategoryId,
			store.CreatedAt,
			store.UpdatedAt,
		)
	}

	newStore, err := repositories.Query(span, dbQuery, readStore)

	return newStore, err
}
