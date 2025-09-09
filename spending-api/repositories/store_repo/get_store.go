package store_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

func (repo *storeRepository) GetStoreById(context context.Context, tx *sql.Tx, id int) (*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetStoreById")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			category_id,
			created_at,
			updated_at
		FROM stores
		WHERE id = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.Query(query, id)
	}

	store, err := repositories.Query(span, dbQuery, readStore)
	return store, err
}

func (repo *storeRepository) GetStoreByUUId(ctx context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetStoreByUUID")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			category_id,
			created_at,
			updated_at
		FROM stores
		WHERE uuid = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.Query(query, uuid)
	}

	store, err := repositories.Query(span, dbQuery, readStore)

	return store, err
}

func (repo *storeRepository) GetStoreByCategoryAndName(ctx context.Context, tx *sql.Tx, categoryId int, name string) (*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetStoreByCategoryAndName")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			category_id,
			created_at,
			updated_at
		FROM stores
		WHERE category_id = $1
		AND name = $2
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.Query(query, categoryId, name)
	}

	store, err := repositories.Query(span, dbQuery, readStore)

	return store, err
}

func (repo *storeRepository) GetStoresByCategoryId(ctx context.Context, tx *sql.Tx, categoryId int) ([]*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetStoresByCategoryId")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			category_id,
			created_at,
			updated_at
		FROM stores
		WHERE category_id = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query, categoryId)
	}

	stores, err := repositories.QueryList(span, dbQuery, readStore)

	return stores, err
}

func (repo *storeRepository) GetStoresByCategoryIds(ctx context.Context, tx *sql.Tx, categoryIds []int) (map[int][]*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetStoresByCategoryIds")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			category_id,
			created_at,
			updated_at
		FROM stores
		WHERE category_id = ANY($1)
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query, pq.Array(categoryIds))
	}

	stores, err := repositories.QueryList(span, dbQuery, readStore)
	if err != nil {
		return nil, err
	}

	storeMap := make(map[int][]*models.Store)
	for _, store := range stores {
		storeMap[store.CategoryId] = append(storeMap[store.CategoryId], store)
	}

	return storeMap, nil
}

func (repo *storeRepository) GetStoreList(ctx context.Context, tx *sql.Tx) ([]*models.Store, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetStoreList")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			category_id,
			created_at,
			updated_at
		FROM stores
		WHERE is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query)
	}

	stores, err := repositories.QueryList(span, dbQuery, readStore)
	return stores, err
}

func readStore(rows *sql.Rows) *models.Store {
	var store models.Store

	err := rows.Scan(
		&store.Id,
		&store.UUId,
		&store.Name,
		&store.CategoryId,
		&store.CreatedAt,
		&store.UpdatedAt)

	utils.CheckError(err)

	return &store
}
