package store_repo

import (
	"context"
	"fmt"
	"spending/models"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *storeRepository) InsertStore(ctx context.Context, store *models.Store) (int, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "DB:InsertStore")
	defer span.End()

	if store == nil {
		return 0, fmt.Errorf("store is nil")
	}

	query := `INSERT INTO stores (
				name,
				category_id,
				created_at,
				updated_at
			) Values ($1, $2, $3, $4)
			 RETURNING id`

	var id int
	err := repo.db.QueryRow(query, store.Name, store.CategoryId, store.CreatedAt, store.UpdatedAt).Scan(&id)

	if err != nil {
		utils.TraceError(span, err)
		return 0, err
	}

	return id, nil
}
