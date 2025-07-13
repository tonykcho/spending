package data_access

import (
	"spending/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDatabase() {
	db := OpenDatabase()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	utils.CheckError(err)
	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	utils.CheckError(err)
	migration.Up()
}

func MigrateTestDatabase() {
	db := OpenDatabase()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	utils.CheckError(err)
	migration, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	utils.CheckError(err)
	migration.Up()
}
