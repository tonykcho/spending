package spending_repo

import (
	"context"
	"database/sql"
	"spending/data_access"
	"spending/models"
	"spending/repositories/category_repo"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func GetSpendingById(context context.Context, id int) *models.SpendingRecord {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingById")
	defer span.End()

	db := data_access.OpenDatabase()

	rows, err := db.Query(`SELECT
							id,
							uuid,
							amount,
							remark,
							spending_date,
							category_id,
							created_at,
							updated_at
						FROM spending_records WHERE id = $1`, id)
	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	record := readSpendingRecord(rows)

	return record
}

func GetSpendingByUUId(context context.Context, uuid uuid.UUID) *models.SpendingRecord {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "GetSpendingByUUId")
	defer span.End()
	db := data_access.OpenDatabase()

	var query string = `SELECT
							id,
							uuid,
							amount,
							remark,
							spending_date,
							category_id,
							created_at,
							updated_at
						FROM spending_records WHERE uuid = $1`

	rows, err := db.Query(query, uuid.String())
	utils.TraceError(span, err)
	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	record := readSpendingRecord(rows)

	return record
}

func GetSpendingList(context context.Context) []*models.SpendingRecord {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()
	db := data_access.OpenDatabase()

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
						ORDER BY spending_date DESC`

	rows, err := db.Query(query)
	utils.TraceError(span, err)
	defer rows.Close()

	var records []*models.SpendingRecord

	for rows.Next() {
		record := readSpendingRecord(rows)
		if record != nil {
			records = append(records, record)
		}
	}

	return records
}

func LoadSpendingCategory(context context.Context, record *models.SpendingRecord) {
	category := category_repo.GetCategoryById(context, record.CategoryId)
	record.Category = category
}

func LoadSpendingListCategory(context context.Context, records []*models.SpendingRecord) {
	categoryIds := make([]int, 0)
	for _, record := range records {
		categoryIds = append(categoryIds, record.CategoryId)
	}

	categories := category_repo.GetCategoryListByIds(context, categoryIds)
	categoryMap := make(map[int]*models.Category)
	for _, category := range categories {
		categoryMap[category.Id] = category
	}

	for _, record := range records {
		record.Category = categoryMap[record.CategoryId]
	}
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
