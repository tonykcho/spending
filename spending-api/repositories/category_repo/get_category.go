package category_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/models"
	"spending/repositories"
	"spending/utils"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *categoryRepository) GetCategoryById(ctx context.Context, tx *sql.Tx, id int) (*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetCategoryById")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query, id)
	}

	category, err := repositories.Query(span, dbQuery, readCategory)

	return category, err
}

func (repo *categoryRepository) GetCategoryByUUId(ctx context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:GetCategoryByUUId")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			created_at,
			updated_at
		FROM categories
		WHERE uuid = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.QueryContext(ctx, query, uuid)
	}

	category, err := repositories.Query(span, dbQuery, readCategory)

	return category, err
}

func (repo *categoryRepository) GetCategoryByName(context context.Context, tx *sql.Tx, name string) (*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryByName")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			name,
			created_at,
			updated_at
		FROM categories
		WHERE name = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query, name)
	}

	category, err := repositories.Query(span, dbQuery, readCategory)

	return category, err
}

func (repo *categoryRepository) GetCategoryList(context context.Context, tx *sql.Tx) ([]*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryList")
	defer span.End()

	var query string = `
		SELECT
			id,
			uuid,
			name,
			created_at,
			updated_at
		FROM categories
		WHERE is_deleted = FALSE
		ORDER BY name
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query)
	}

	categories, err := repositories.QueryList(span, dbQuery, readCategory)

	return categories, err
}

func (repo *categoryRepository) GetCategoryListByIds(context context.Context, tx *sql.Tx, ids []int) ([]*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryListByIds")
	defer span.End()

	if len(ids) == 0 {
		return []*models.Category{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`SELECT
							id,
							uuid,
							name,
							created_at,
							updated_at
						  FROM categories
						  WHERE id IN (%s)
						  AND is_deleted = FALSE
						  ORDER BY id`, strings.Join(placeholders, ","))

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query, args...)
	}

	categories, err := repositories.QueryList(span, dbQuery, readCategory)

	return categories, err
}

func (repo *categoryRepository) LoadStoresForCategory(ctx context.Context, tx *sql.Tx, category *models.Category) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:LoadStoresForCategory")
	defer span.End()

	if category == nil {
		return nil
	}

	stores, err := repo.store_repo.GetStoresByCategoryId(ctx, tx, category.Id)
	if err != nil {
		return err
	}

	category.Stores = stores

	return nil
}

func (repo *categoryRepository) LoadStoresForCategories(ctx context.Context, tx *sql.Tx, categories []*models.Category) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:LoadStoresForCategories")
	defer span.End()

	if len(categories) == 0 {
		return nil
	}

	categoryIds := make([]int, len(categories))
	for i, category := range categories {
		categoryIds[i] = category.Id
	}

	storesMap, err := repo.store_repo.GetStoresByCategoryIds(ctx, tx, categoryIds)
	if err != nil {
		return err
	}

	for _, category := range categories {
		stores, ok := storesMap[category.Id]
		if ok {
			category.Stores = stores
		} else {
			category.Stores = []*models.Store{}
		}
	}

	return nil
}

func readCategory(rows *sql.Rows) *models.Category {
	var category models.Category

	err := rows.Scan(
		&category.Id,
		&category.UUId,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt)

	utils.CheckError(err)
	return &category
}
