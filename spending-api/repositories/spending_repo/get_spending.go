package spending_repo

import (
	"context"
	"database/sql"
	"spending/models"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func (repo *spendingRepository) GetSpendingById(context context.Context, id int) (*models.SpendingRecord, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingById")
	defer span.End()

	rows, err := repo.db.Query(`SELECT
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
						AND is_deleted = FALSE`, id)
	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

	record := readSpendingRecord(rows)

	return record, nil
}

func (repo *spendingRepository) GetSpendingByUUId(context context.Context, uuid uuid.UUID) (*models.SpendingRecord, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "GetSpendingByUUId")
	defer span.End()

	var query string = `SELECT
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
						AND is_deleted = FALSE`

	rows, err := repo.db.Query(query, uuid.String())
	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

	record := readSpendingRecord(rows)

	return record, nil
}

func (repo *spendingRepository) GetSpendingList(context context.Context) ([]*models.SpendingRecord, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()

	var query string = `SELECT 
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
						ORDER BY spending_date DESC`

	rows, err := repo.db.Query(query)
	utils.TraceError(span, err)
	defer rows.Close()

	var records []*models.SpendingRecord

	for rows.Next() {
		record := readSpendingRecord(rows)
		if record != nil {
			records = append(records, record)
		}
	}

	return records, err
}

func (repo *spendingRepository) LoadSpendingCategory(context context.Context, record *models.SpendingRecord) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()

	category, err := repo.category_repo.GetCategoryById(context, record.CategoryId)
	if err != nil {
		utils.TraceError(span, err)
		return err
	}

	record.Category = category
	return nil
}

func (repo *spendingRepository) LoadSpendingListCategory(context context.Context, records []*models.SpendingRecord) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()

	categoryIds := make([]int, 0)
	for _, record := range records {
		categoryIds = append(categoryIds, record.CategoryId)
	}

	categories, err := repo.category_repo.GetCategoryListByIds(context, categoryIds)
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
