package zintegration_test

import (
	"context"
	"database/sql"
	"os"
	"spending/data_access"
	"spending/utils"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	log.Info().Msg("Starting PostgreSQL test container...")
	pgContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	utils.CheckError(err)
	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	log.Info().Msg(connectionString)
	utils.CheckError(err)
	data_access.DB, err = sql.Open("postgres", connectionString)
	utils.CheckError(err)

	// Wait for the database to be ready
	err = data_access.DB.Ping()
	for i := 0; i < 10 && err != nil; i++ {
		log.Info().Msg("Waiting for database to be ready...")
		time.Sleep(1 * time.Second)
		err = data_access.DB.Ping()
	}

	log.Info().Msg("Start Migrating test database...")
	data_access.MigrateTestDatabase()

	code := m.Run()
	pgContainer.Terminate(ctx)
	os.Exit(code)
}
