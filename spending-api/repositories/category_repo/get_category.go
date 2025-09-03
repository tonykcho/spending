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

func GetCategoryById(context context.Context, id int) *models.Category {
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
						FROM categories WHERE id = $1`, id)

	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	category := readCategory(rows)

	return category
}

func GetCategoryByUUId(context context.Context, uuid uuid.UUID) *models.Category {
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
						FROM categories WHERE uuid = $1`, uuid)

	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	category := readCategory(rows)

	return category
}

func GetCategoryList(context context.Context) []*models.Category {
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
						FROM categories`

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

	return categories
}

func GetCategoryListByIds(context context.Context, ids []int) []*models.Category {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryListByIds")
	defer span.End()

	if len(ids) == 0 {
		return []*models.Category{}
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
						  FROM categories WHERE id IN (%s)`, strings.Join(placeholders, ","))
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

	return categories
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
