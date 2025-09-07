package store_handlers

import (
	"database/sql"
	"net/http"
	"spending/repositories"
	"spending/repositories/store_repo"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type deleteStoreHandler struct {
	store_repo   store_repo.StoreRepository
	unit_of_work repositories.UnitOfWork
}

func NewDeleteStoreHandler(storeRepo store_repo.StoreRepository, unitOfWork repositories.UnitOfWork) *deleteStoreHandler {
	return &deleteStoreHandler{
		store_repo:   storeRepo,
		unit_of_work: unitOfWork,
	}
}

func (handler *deleteStoreHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "DeleteStoreHandler")
	defer span.End()

	routerVars := mux.Vars(request)
	storeUUId, err := uuid.Parse(routerVars["id"])

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	status := http.StatusInternalServerError

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		txErr := handler.store_repo.DeleteStore(ctx, tx, storeUUId)
		if txErr != nil {
			status = http.StatusInternalServerError
			return txErr
		}
		return nil
	})

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), status)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
