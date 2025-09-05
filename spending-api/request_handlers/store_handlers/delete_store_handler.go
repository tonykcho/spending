package store_handlers

import (
	"net/http"
	"spending/repositories/store_repo"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type deleteStoreHandler struct {
	store_repo store_repo.StoreRepository
}

func NewDeleteStoreHandler(storeRepo store_repo.StoreRepository) *deleteStoreHandler {
	return &deleteStoreHandler{
		store_repo: storeRepo,
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

	err = handler.store_repo.DeleteStore(ctx, storeUUId)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
