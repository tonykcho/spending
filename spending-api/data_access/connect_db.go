package data_access

import (
	"database/sql"
	"spending/utils"
)

func OpenDatabase() *sql.DB {
	connectionString := utils.GetDatabaseConnection()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		utils.CheckError(err)
	}
	return db
}

func OpenPostgres() *sql.DB {
	connectionString := utils.GetPostgresConnection()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		utils.CheckError(err)
	}
	return db
}
