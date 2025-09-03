package category_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/data_access"
	"spending/models"
	"spending/utils"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func GetCategoryById(context context.Context, id int) (*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryById")
	defer span.End()

	db := data_access.OpenDatabase()

	rows, err := db.Query(`SELECT
							id,
							uuid,
							name,
							created_at,
							updated_at
						FROM categories
						WHERE id = $1
						AND is_deleted = FALSE`, id)

	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

	category := readCategory(rows)

	return category, err
}

func GetCategoryByUUId(context context.Context, uuid uuid.UUID) (*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryByUUId")
	defer span.End()

	db := data_access.OpenDatabase()

	rows, err := db.Query(`SELECT
							id,
							uuid,
							name,
							created_at,
							updated_at
						FROM categories
						WHERE uuid = $1
						AND is_deleted = FALSE`, uuid)

	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

	category := readCategory(rows)

	return category, err
}

func GetCategoryByName(context context.Context, name string) (*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryByName")
	defer span.End()

	db := data_access.OpenDatabase()

	rows, err := db.Query(`SELECT
							id,
							uuid,
							name,
							created_at,
							updated_at
						FROM categories
						WHERE name = $1
						AND is_deleted = FALSE`, name)

	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

	category := readCategory(rows)

	return category, err
}

func GetCategoryList(context context.Context) ([]*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryList")
	defer span.End()

	db := data_access.OpenDatabase()

	var query string = `SELECT
							id,
							uuid,
							name,
							created_at,
							updated_at
						FROM categories
						WHERE is_deleted = FALSE`

	rows, err := db.Query(query)
	utils.TraceError(span, err)
	defer rows.Close()

	var categories []*models.Category = make([]*models.Category, 0)

	for rows.Next() {
		category := readCategory(rows)
		if category != nil {
			categories = append(categories, category)
		}
	}

	return categories, err
}

func GetCategoryListByIds(context context.Context, ids []int) ([]*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryListByIds")
	defer span.End()

	if len(ids) == 0 {
		return []*models.Category{}, nil
	}

	db := data_access.OpenDatabase()

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
						  AND is_deleted = FALSE`, strings.Join(placeholders, ","))
	rows, err := db.Query(query, args...)

	utils.TraceError(span, err)
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		category := readCategory(rows)
		if category != nil {
			categories = append(categories, category)
		}
	}

	return categories, err
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
