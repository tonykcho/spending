package store_handlers

import (
	"net/http"
	"spending/repositories/store_repo"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type getStoreHandler struct {
	store_repo store_repo.StoreRepository
}

func NewGetStoreHandler(storeRepo store_repo.StoreRepository) *getStoreHandler {
	return &getStoreHandler{
		store_repo: storeRepo,
	}
}

func (handler *getStoreHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "GetStoreHandler")
	defer span.End()

	routerVars := mux.Vars(request)
	storeUUId, err := uuid.Parse(routerVars["id"])

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	store, err := handler.store_repo.GetStoreByUUId(ctx, storeUUId)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if store == nil {
		http.Error(writer, "Record not found", http.StatusNotFound)
		return
	}

	err = utils.Encode(ctx, writer, http.StatusOK, store)
	utils.TraceError(span, err)
}
