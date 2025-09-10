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

func (repo *storeRepository) InsertStores(ctx context.Context, tx *sql.Tx, stores []*models.Store) ([]*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:InsertStores")
	defer span.End()

	if len(stores) == 0 {
		return []*models.Store{}, nil
	}

	query := `INSERT INTO stores (
				name,
				category_id,
				created_at,
				updated_at
			) Values `

	params := []interface{}{}
	for i, store := range stores {
		paramIdx := i*4 + 1
		query += fmt.Sprintf("($%d, $%d, $%d, $%d),", paramIdx, paramIdx+1, paramIdx+2, paramIdx+3)
		params = append(params, store.Name, store.CategoryId, store.CreatedAt, store.UpdatedAt)
	}

	query = query[:len(query)-1] // Remove trailing comma
	query += `
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
		return dbTx.QueryContext(ctx, query, params...)
	}

	newStores, err := repositories.QueryList(span, dbQuery, readStore)

	return newStores, err
}
