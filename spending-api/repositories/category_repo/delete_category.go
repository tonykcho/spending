package category_repo

import (
	"context"
	"spending/data_access"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func DeleteCategory(context context.Context, uuid uuid.UUID) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:DeleteCategory")
	defer span.End()

	db := data_access.OpenDatabase()

	_, err := db.Exec(`UPDATE categories SET is_deleted = TRUE, deleted_at = NOW() WHERE uuid = $1`, uuid)
	utils.TraceError(span, err)
	return err
}
