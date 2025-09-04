package data_access

import (
	"database/sql"
	"spending/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDatabase(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	utils.CheckError(err)
	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	utils.CheckError(err)
	err = migration.Up()

	if err != nil && err != migrate.ErrNoChange {
		utils.CheckError(err)
	}
}

func MigrateTestDatabase(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	utils.CheckError(err)
	migration, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	utils.CheckError(err)
	migration.Up()
}
