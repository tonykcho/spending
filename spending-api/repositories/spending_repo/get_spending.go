package spending_repo

import (
	"context"
	"database/sql"
	"fmt"
	"spending/data_access"
	"spending/models"
	"spending/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func GetSpendingById(context context.Context, id int) *models.SpendingRecord {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingById")
	defer span.End()

	db := data_access.OpenDatabase()

	var queryTemplate string = `SELECT * FROM spending_records WHERE id = %d`
	var query = fmt.Sprintf(queryTemplate, id)

	rows, err := db.Query(query)
	utils.CheckError(err)
	defer rows.Close()

	record := readSpendingRecord(rows)

	return record
}

func GetSpendingByUUId(context context.Context, uid uuid.UUID) *models.SpendingRecord {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "GetSpendingByUUId")
	defer span.End()
	db := data_access.OpenDatabase()

	var queryTemplate string = "SELECT * FROM spending_records WHERE uuid = '%s'"
	var query string = fmt.Sprintf(queryTemplate, uid.String())

	rows, err := db.Query(query)
	utils.CheckError(err)
	defer rows.Close()

	record := readSpendingRecord(rows)

	return record
}

func GetSpendingList(context context.Context) []*models.SpendingRecord {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DB:GetSpendingList")
	defer span.End()
	db := data_access.OpenDatabase()

	var query string = "SELECT * FROM spending_records ORDER BY spending_date DESC"

	rows, err := db.Query(query)
	utils.CheckError(err)
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

func readSpendingRecord(rows *sql.Rows) *models.SpendingRecord {
	if !rows.Next() {
		return nil
	}

	var record models.SpendingRecord

	err := rows.Scan(
		&record.Id,
		&record.UUId,
		&record.Amount,
		&record.Remark,
		&record.SpendingDate,
		&record.Category,
		&record.CreatedAt,
		&record.UpdatedAt)

	utils.CheckError(err)
	return &record
}
