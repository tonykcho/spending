package store_repo

import (
	"context"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *storeRepository) DeleteStore(ctx context.Context, uuid uuid.UUID) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:DeleteStore")
	defer span.End()

	query := `
		UPDATE stores 
		SET is_deleted = TRUE, deleted_at = NOW() 
		WHERE uuid = $1
	`

	_, err := repo.db.Exec(query, uuid)
	utils.TraceError(span, err)
	return err
}
