package category_repo

import (
	"context"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *categoryRepository) DeleteCategory(context context.Context, uuid uuid.UUID) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:DeleteCategory")
	defer span.End()

	_, err := repo.db.Exec(`UPDATE categories SET is_deleted = TRUE, deleted_at = NOW() WHERE uuid = $1`, uuid)
	utils.TraceError(span, err)
	return err
}
