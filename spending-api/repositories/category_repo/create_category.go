package category_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *categoryRepository) InsertCategory(context context.Context, tx *sql.Tx, category models.Category) (*models.Category, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:InsertCategory")
	defer span.End()

	query := `
	INSERT INTO categories (
		name,
		created_at,
		updated_at
	) Values ($1, $2, $3)
		RETURNING
			id,
			uuid,
			name,
			created_at,
			updated_at
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	var dbQuery = func() (*sql.Rows, error) {
		return dbTx.QueryContext(context, query)
	}

	newCategory, err := repositories.Query(span, dbQuery, readCategory)

	utils.TraceError(span, err)
	return newCategory, err
}
