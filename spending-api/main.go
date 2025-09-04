package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"spending/data_access"
	"spending/repositories/category_repo"
	"spending/repositories/spending_repo"
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

	"github.com/rs/cors"
)

type Container struct {
	CategoryRepository category_repo.CategoryRepository
	SpendingRepository spending_repo.SpendingRepository

	CreateCategoryHandler  category_handlers.CreateCategoryHandler
	DeleteCategoryHandler  category_handlers.DeleteCategoryHandler
	GetCategoryHandler     category_handlers.GetCategoryHandler
	GetCategoryListHandler category_handlers.GetCategoryListHandler
	UpdateCategoryHandler  category_handlers.UpdateCategoryHandler

	CreateSpendingHandler  spending_handlers.CreateSpendingHandler
	GetSpendingHandler     spending_handlers.GetSpendingHandler
	GetSpendingListHandler spending_handlers.GetSpendingListHandler
	DeleteSpendingHandler  spending_handlers.DeleteSpendingHandler
}

func NewContainer(db *sql.DB) *Container {
	categoryRepo := category_repo.NewCategoryRepository(db)
	spendingRepo := spending_repo.NewSpendingRepository(db)

	return &Container{
		CategoryRepository:     categoryRepo,
		SpendingRepository:     spendingRepo,
		CreateCategoryHandler:  category_handlers.NewCreateCategoryHandler(categoryRepo),
		DeleteCategoryHandler:  category_handlers.NewDeleteCategoryHandler(categoryRepo),
		GetCategoryHandler:     category_handlers.NewGetCategoryHandler(categoryRepo),
		GetCategoryListHandler: category_handlers.NewGetCategoryListHandler(categoryRepo),
		UpdateCategoryHandler:  category_handlers.NewUpdateCategoryHandler(categoryRepo),

		CreateSpendingHandler:  spending_handlers.NewCreateSpendingHandler(spendingRepo, categoryRepo),
		GetSpendingHandler:     spending_handlers.NewGetSpendingHandler(spendingRepo),
		GetSpendingListHandler: spending_handlers.NewGetSpendingListHandler(spendingRepo),
		DeleteSpendingHandler:  spending_handlers.NewDeleteSpendingHandler(spendingRepo),
	}
}

func main() {
	// Make sure database is up.
	data_access.CreateDatabase()

	router := mux.NewRouter()
	db := data_access.OpenDatabase()
	defer db.Close()

	configureLogging()
	data_access.MigrateDatabase(db)
	configureEndpoints(router, db)
	configureOpenTelemetry()

	handler := cors.AllowAll().Handler(router)

	log.Info().Msg("Server is listening on port 8001")
	http.ListenAndServe(":8001", handler)
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

func configureEndpoints(router *mux.Router, db *sql.DB) {
	container := NewContainer(db)

	router.HandleFunc("/spending/{id}", container.GetSpendingHandler.Handle).Methods("GET")
	router.HandleFunc("/spending", container.GetSpendingListHandler.Handle).Methods("GET")
	router.HandleFunc("/spending", container.CreateSpendingHandler.Handle).Methods("POST")
	router.HandleFunc("/spending/{id}", container.DeleteSpendingHandler.Handle).Methods("DELETE")

	router.HandleFunc("/categories/{id}", container.GetCategoryHandler.Handle).Methods("GET")
	router.HandleFunc("/categories", container.GetCategoryListHandler.Handle).Methods("GET")
	router.HandleFunc("/categories", container.CreateCategoryHandler.Handle).Methods("POST")
	router.HandleFunc("/categories/{id}", container.UpdateCategoryHandler.Handle).Methods("PUT")
	router.HandleFunc("/categories/{id}", container.DeleteCategoryHandler.Handle).Methods("DELETE")

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
