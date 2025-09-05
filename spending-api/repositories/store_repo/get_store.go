package store_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *storeRepository) GetStoreById(context context.Context, id int) (*models.Store, error) {
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

	var dbQuery = func() (*sql.Rows, error) {
		return repo.db.Query(query, id)
	}

	store, err := repositories.Query(span, dbQuery, readStore)
	return store, err
}

func (repo *storeRepository) GetStoreByUUId(ctx context.Context, uuid uuid.UUID) (*models.Store, error) {
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

	var dbQuery = func() (*sql.Rows, error) {
		return repo.db.Query(query, uuid)
	}

	store, err := repositories.Query(span, dbQuery, readStore)

	return store, err
}

func (repo *storeRepository) GetStoreByCategoryAndName(ctx context.Context, categoryId int, name string) (*models.Store, error) {
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

	var dbQuery = func() (*sql.Rows, error) {
		return repo.db.Query(query, categoryId, name)
	}

	store, err := repositories.Query(span, dbQuery, readStore)

	return store, err
}

func (repo *storeRepository) GetStoreList(ctx context.Context) ([]*models.Store, error) {
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

	dbQuery := func() (*sql.Rows, error) {
		return repo.db.Query(query)
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
