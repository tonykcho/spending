package data_access

import (
	"database/sql"

	"spending/utils"

	"github.com/rs/zerolog/log"
)

func CreateDatabase() {
	db := OpenPostgres()

	if isDatabaseExist(db) {
		return
	}
	query := `CREATE DATABASE spending`

	_, err := db.Exec(query)
	utils.CheckError(err)

	log.Info().Msg("Database 'spending' created successfully")
}

func isDatabaseExist(db *sql.DB) bool {
	query := `Select 1 FROM pg_database WHERE datname = 'spending'`

	rows := db.QueryRow(query)

	var result int
	err := rows.Scan(&result)

	var exist bool = false

	switch err {
	case nil:
		log.Info().Msg("Database already exists")
		exist = true
	case sql.ErrNoRows:
		log.Info().Msg("Database does not exist")
		exist = false
	default:
		utils.CheckError(err)
	}

	return exist
}
