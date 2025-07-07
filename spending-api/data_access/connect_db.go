package data_access

import (
	"database/sql"
	"spending/utils"
)

func OpenDatabase() *sql.DB {
	connectionString := "user=postgres password=pwd dbname=spending sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		utils.CheckError(err)
	}
	return db
}

func OpenPostgres() *sql.DB {
	connectionString := "user=postgres password=pwd sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		utils.CheckError(err)
	}
	return db
}
