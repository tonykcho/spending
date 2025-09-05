package store_handlers

import (
	"net/http"
	"spending/mappers"
	"spending/repositories/store_repo"
	"spending/request_handlers"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

type getStoreListHandler struct {
	store_repo store_repo.StoreRepository
}

func NewGetStoreListHandler(storeRepo store_repo.StoreRepository) request_handlers.RequestHandler {
	return &getStoreListHandler{
		store_repo: storeRepo,
	}
}

func (handler *getStoreListHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "GetStoreListHandler")
	defer span.End()

	stores, err := handler.store_repo.GetStoreList(context)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response := mappers.MapStoreList(stores)

	err = utils.Encode(context, writer, http.StatusOK, response)
	utils.TraceError(span, err)
}
