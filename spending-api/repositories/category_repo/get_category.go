package category_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/data_access"
	"spending/models"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func GetCategoryById(context context.Context, id int) *models.Category {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryById")
	defer span.End()

	db := data_access.OpenDatabase()

	var queryTemplate string = `SELECT * FROM categories WHERE id = %d`
	var query = fmt.Sprintf(queryTemplate, id)

	rows, err := db.Query(query)
	utils.TraceError(span, err)
	defer rows.Close()

	category := readCategory(rows)

	return category
}

func GetCategoryByUUId(context context.Context, uuid uuid.UUID) *models.Category {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryByUUId")
	defer span.End()

	db := data_access.OpenDatabase()

	var queryTemplate string = `SELECT * FROM categories WHERE uuid = '%s'`
	var query = fmt.Sprintf(queryTemplate, uuid)

	rows, err := db.Query(query)
	utils.TraceError(span, err)
	defer rows.Close()

	category := readCategory(rows)

	return category
}

func GetCategoryList(context context.Context) []*models.Category {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetCategoryList")
	defer span.End()

	db := data_access.OpenDatabase()

	var query string = "SELECT * FROM categories"

	rows, err := db.Query(query)
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
	if !rows.Next() {
		return nil
	}

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
