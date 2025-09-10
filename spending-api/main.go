package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"spending/data_access"
	"spending/repositories"
	"spending/repositories/category_repo"
	"spending/repositories/spending_repo"
	"spending/repositories/store_repo"
	"spending/request_handlers"
	"spending/request_handlers/category_handlers"
	"spending/request_handlers/spending_handlers"
	"spending/request_handlers/store_handlers"
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
	StoreRepository    store_repo.StoreRepository
	UnitOfWork         repositories.UnitOfWork

	CreateCategoryHandler  request_handlers.RequestHandler
	DeleteCategoryHandler  request_handlers.RequestHandler
	GetCategoryHandler     request_handlers.RequestHandler
	GetCategoryListHandler request_handlers.RequestHandler
	UpdateCategoryHandler  request_handlers.RequestHandler

	CreateSpendingHandler  request_handlers.RequestHandler
	GetSpendingHandler     request_handlers.RequestHandler
	GetSpendingListHandler request_handlers.RequestHandler
	DeleteSpendingHandler  request_handlers.RequestHandler

	CreateStoreHandler  request_handlers.RequestHandler
	DeleteStoreHandler  request_handlers.RequestHandler
	GetStoreHandler     request_handlers.RequestHandler
	GetStoreListHandler request_handlers.RequestHandler
}

func NewContainer(db *sql.DB) *Container {
	storeRepo := store_repo.NewStoreRepository(db)
	categoryRepo := category_repo.NewCategoryRepository(db, storeRepo)
	spendingRepo := spending_repo.NewSpendingRepository(db, categoryRepo)
	unitOfWork := repositories.NewUnitOfWork(db)

	return &Container{
		CategoryRepository: categoryRepo,
		SpendingRepository: spendingRepo,
		StoreRepository:    storeRepo,
		UnitOfWork:         unitOfWork,

		CreateCategoryHandler:  category_handlers.NewCreateCategoryHandler(categoryRepo, storeRepo, unitOfWork),
		DeleteCategoryHandler:  category_handlers.NewDeleteCategoryHandler(categoryRepo, unitOfWork),
		GetCategoryHandler:     category_handlers.NewGetCategoryHandler(categoryRepo),
		GetCategoryListHandler: category_handlers.NewGetCategoryListHandler(categoryRepo),
		UpdateCategoryHandler:  category_handlers.NewUpdateCategoryHandler(categoryRepo, unitOfWork),

		CreateSpendingHandler:  spending_handlers.NewCreateSpendingHandler(spendingRepo, categoryRepo, unitOfWork),
		GetSpendingHandler:     spending_handlers.NewGetSpendingHandler(spendingRepo),
		GetSpendingListHandler: spending_handlers.NewGetSpendingListHandler(spendingRepo),
		DeleteSpendingHandler:  spending_handlers.NewDeleteSpendingHandler(spendingRepo, unitOfWork),

		CreateStoreHandler:  store_handlers.NewCreateStoreHandler(storeRepo, categoryRepo, unitOfWork),
		DeleteStoreHandler:  store_handlers.NewDeleteStoreHandler(storeRepo, unitOfWork),
		GetStoreHandler:     store_handlers.NewGetStoreHandler(storeRepo),
		GetStoreListHandler: store_handlers.NewGetStoreListHandler(storeRepo),
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

	router.HandleFunc("/stores/{id}", container.GetStoreHandler.Handle).Methods("GET")
	router.HandleFunc("/stores", container.GetStoreListHandler.Handle).Methods("GET")
	router.HandleFunc("/stores", container.CreateStoreHandler.Handle).Methods("POST")
	router.HandleFunc("/stores/{id}", container.DeleteStoreHandler.Handle).Methods("DELETE")

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
