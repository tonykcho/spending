package spending_repo

import (
	"spending/data_access"
	"spending/models"
	"spending/utils"
)

func InsertSpendingRecord(record models.SpendingRecord) int {
	db := data_access.OpenDatabase()
	defer db.Close()

	// Create query to insert a new spending record
	query := `INSERT INTO spending_records (
				amount,
				remark,
				spending_date,
				category,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6)
			 RETURNING id`

	var id int
	err := db.QueryRow(query,
		record.Amount,
		record.Remark,
		record.SpendingDate,
		record.Category,
		record.CreatedAt,
		record.UpdatedAt,
	).Scan(&id)

	utils.CheckError(err)
	return id
}
