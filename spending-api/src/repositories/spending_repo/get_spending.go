package spending_repo

import (
	"database/sql"
	"fmt"
	"spending/data_access"
	"spending/models"
	"spending/utils"

	"github.com/google/uuid"
)

func GetSpendingById(id int) *models.SpendingRecord {
	db := data_access.OpenDatabase()
	defer db.Close()

	var queryTemplate string = `SELECT * FROM spending WHERE spending.id = %d`
	var query = fmt.Sprintf(queryTemplate, id)

	rows, err := db.Query(query)
	utils.CheckError(err)
	defer rows.Close()

	record := readSpendingRecord(rows)

	return record
}

func GetSpendingByUUId(uid uuid.UUID) *models.SpendingRecord {
	db := data_access.OpenDatabase()
	defer db.Close()

	var queryTemplate string = "SELECT * FROM spending_records WHERE spending_records.uuid = %s"
	var query string = fmt.Sprintf(queryTemplate, uid.String())

	rows, err := db.Query(query)
	utils.CheckError(err)
	defer rows.Close()

	record := readSpendingRecord(rows)

	return record
}

func readSpendingRecord(rows *sql.Rows) *models.SpendingRecord {
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

	if err == sql.ErrNoRows {
		return nil
	} else {
		utils.CheckError(err)
	}

	return &record
}
