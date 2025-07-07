package main

import (
	"net/http"
	"os"
	"spending/data_access"
	"spending/request_handlers/spending_handlers"
	"spending/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	router := mux.NewRouter()

	configureLogging()
	configureDatabase()
	configureEndpoints(router)

	log.Info().Msg("Server is listening on port 8080")
	http.ListenAndServe(":8001", router)
}

func configureLogging() {
	dateStr := time.Now().Format("2006-01-02")
	logDir := "logs"
	logFilePath := logDir + "/spending-api-" + dateStr + ".log"

	// Ensure the log directory exists, 0755 stands for read/write/execute permissions for owner, and read/execute for group and others
	err := os.MkdirAll(logDir, 0755)
	utils.CheckError(err)

	// Open the log file for appending, create if not exists
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.CheckError(err)

	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}

func configureDatabase() {
	data_access.CreateDatabase()
	data_access.MigrateDatabase()
}

func configureEndpoints(router *mux.Router) {
	log.Info().Msg("Adding spending endpoints")
	spending_handlers.RegisterSpendingHandlers(router)
}
