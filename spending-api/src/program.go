package main

import (
	"net/http"
	"spending/data_access"

	"github.com/rs/zerolog/log"
)

func main() {
	initialize()

	log.Info().Msg("Server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func initialize() {
	configureDatabase()
}

func configureDatabase() {
	data_access.CreateDatabase()
	data_access.MigrateDatabase()
}
