package main

import (
	"net/http"
	"spending/data_access"
	"spending/request_handlers/spending_handlers"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	router := mux.NewRouter()

	configureDatabase()
	configureEndpoints(router)

	log.Info().Msg("Server is listening on port 8080")
	http.ListenAndServe(":8001", router)
}

func configureDatabase() {
	data_access.CreateDatabase()
	data_access.MigrateDatabase()
}

func configureEndpoints(router *mux.Router) {
	log.Info().Msg("Adding spending endpoints")
	spending_handlers.RegisterSpendingHandlers(router)
}
