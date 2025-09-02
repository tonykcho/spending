package category_repo

import (
	"context"
	"spending/data_access"
	"spending/models"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func InsertCategory(context context.Context, category models.Category) int {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:InsertCategory")
	defer span.End()

	db := data_access.OpenDatabase()

	query := `INSERT INTO categories (
				name,
				created_at,
				updated_at
			) Values ($1, $2, $3)
			 RETURNING id`

	var id int
	err := db.QueryRow(query, category.Name, category.CreatedAt, category.UpdatedAt).Scan(&id)
	utils.TraceError(span, err)
	return id
}
