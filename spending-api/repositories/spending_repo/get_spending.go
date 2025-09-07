package spending_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/repositories"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *spendingRepository) GetSpendingById(context context.Context, tx *sql.Tx, id int) (*models.SpendingRecord, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingById")
	defer span.End()

	query := `
		SELECT
			id,
			uuid,
			amount,
			remark,
			spending_date,
			category_id,
			created_at,
			updated_at
		FROM spending_records
		WHERE id = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query, id)
	}

	record, err := repositories.Query(span, dbQuery, readSpendingRecord)

	return record, err
}

func (repo *spendingRepository) GetSpendingByUUId(context context.Context, tx *sql.Tx, uuid uuid.UUID) (*models.SpendingRecord, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "GetSpendingByUUId")
	defer span.End()

	var query string = `
		SELECT
			id,
			uuid,
			amount,
			remark,
			spending_date,
			category_id,
			created_at,
			updated_at
		FROM spending_records
		WHERE uuid = $1
		AND is_deleted = FALSE
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query, uuid)
	}

	record, err := repositories.Query(span, dbQuery, readSpendingRecord)

	return record, err
}

func (repo *spendingRepository) GetSpendingList(context context.Context, tx *sql.Tx) ([]*models.SpendingRecord, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()

	var query string = `
		SELECT 
			id,
			uuid,
			amount,
			remark,
			spending_date,
			category_id,
			created_at,
			updated_at
		FROM spending_records
		WHERE is_deleted = FALSE
		ORDER BY spending_date DESC
	`

	var dbTx repositories.DbTx = repo.db
	if tx != nil {
		dbTx = tx
	}

	dbQuery := func() (*sql.Rows, error) {
		return dbTx.Query(query)
	}

	records, err := repositories.QueryList(span, dbQuery, readSpendingRecord)

	return records, err
}

func (repo *spendingRepository) LoadSpendingCategory(context context.Context, tx *sql.Tx, record *models.SpendingRecord) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()

	category, err := repo.category_repo.GetCategoryById(context, tx, record.CategoryId)
	if err != nil {
		utils.TraceError(span, err)
		return err
	}

	record.Category = category
	return nil
}

func (repo *spendingRepository) LoadSpendingListCategory(context context.Context, tx *sql.Tx, records []*models.SpendingRecord) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()

	categoryIds := make([]int, 0)
	for _, record := range records {
		categoryIds = append(categoryIds, record.CategoryId)
	}

	categories, err := repo.category_repo.GetCategoryListByIds(context, tx, categoryIds)
	if err != nil {
		utils.TraceError(span, err)
		return err
	}

	categoryMap := make(map[int]*models.Category)
	for _, category := range categories {
		categoryMap[category.Id] = category
	}

	for _, record := range records {
		record.Category = categoryMap[record.CategoryId]
	}

	return nil
}

func readSpendingRecord(rows *sql.Rows) *models.SpendingRecord {
	var record models.SpendingRecord

	err := rows.Scan(
		&record.Id,
		&record.UUId,
		&record.Amount,
		&record.Remark,
		&record.SpendingDate,
		&record.CategoryId,
		&record.CreatedAt,
		&record.UpdatedAt)

	utils.CheckError(err)
	return &record
}
