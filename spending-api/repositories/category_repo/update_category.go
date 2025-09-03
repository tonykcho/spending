package category_repo

import (
	"context"
	"spending/models"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *categoryRepository) UpdateCategory(context context.Context, category *models.Category) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:UpdateCategory")
	defer span.End()

	_, err := repo.db.Exec(`UPDATE categories SET
							name = $1,
							updated_at = $2
						WHERE id = $3`, category.Name, category.UpdatedAt, category.Id)

	utils.TraceError(span, err)
	return err
}
