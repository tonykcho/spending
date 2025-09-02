package main

import (
	"context"
	"net/http"
	"os"
	"spending/data_access"
	"spending/request_handlers/category_handlers"
	"spending/request_handlers/spending_handlers"
	"spending/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
	router := mux.NewRouter()

	configureLogging()
	configureDatabase()
	configureEndpoints(router)
	configureOpenTelemetry()

	log.Info().Msg("Server is listening on port 8001")
	http.ListenAndServe(":8001", router)
	data_access.CloseDatabase() // Ensure the database connection is closed when the server stops
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
	router.HandleFunc("/spending", spending_handlers.GetSpendingListHandler).Methods("GET")
	router.HandleFunc("/spending/{id}", spending_handlers.GetSpendingRequestHandler).Methods("GET")
	router.HandleFunc("/spending", spending_handlers.CreateSpendingRequestHandler).Methods("POST")

	router.HandleFunc("/categories", category_handlers.CreateCategoryRequestHandler).Methods("POST")

	router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
}

func configureOpenTelemetry() {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)

	utils.CheckError(err)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("spending-api"),
		)),
	)

	otel.SetTracerProvider(tp)
	log.Info().Msg("OpenTelemetry configured with OTLP exporter")
}
