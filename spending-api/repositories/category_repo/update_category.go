package category_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func (repo *categoryRepository) UpdateCategory(context context.Context, tx *sql.Tx, category *models.Category) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:UpdateCategory")
	defer span.End()

	query := `
		UPDATE categories SET
			name = $1,
			updated_at = $2
		WHERE id = $3
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	_, err := dbTx.ExecContext(context, query, category.Name, category.UpdatedAt, category.Id)

	utils.TraceError(span, err)
	return err
}
