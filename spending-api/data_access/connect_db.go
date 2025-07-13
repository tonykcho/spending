package data_access

import (
	"database/sql"
	"spending/utils"
)

var DB *sql.DB // Singleton database connection

func OpenDatabase() *sql.DB {
	if DB != nil {
		return DB // Return existing connection if already opened
	}

	connectionString := "user=postgres password=pwd dbname=spending sslmode=disable"
	DB, err := sql.Open("postgres", connectionString)
	if err != nil {
		utils.CheckError(err)
	}
	return DB
}

func OpenPostgres() *sql.DB {
	connectionString := "user=postgres password=pwd sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		utils.CheckError(err)
	}
	return db
}

func CloseDatabase() {
	if DB == nil {
		return // No connection to close
	}
	err := DB.Close()
	utils.CheckError(err)
	DB = nil // Reset the DB variable to nil after closing
}
