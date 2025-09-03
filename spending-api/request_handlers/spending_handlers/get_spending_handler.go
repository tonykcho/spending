package spending_handlers

import (
	"encoding/json"
	"net/http"
	"spending/mappers"
	"spending/repositories/spending_repo"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type GetSpendingHandler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
}

type getSpendingHandler struct {
	spending_repo spending_repo.SpendingRepository
}

func NewGetSpendingHandler(spendingRepo spending_repo.SpendingRepository) GetSpendingHandler {
	return &getSpendingHandler{
		spending_repo: spendingRepo,
	}
}

func (handler *getSpendingHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "DB:GetSpendingByUUId")
	defer span.End()

	routerVars := mux.Vars(request)

	spendingUUId, err := uuid.Parse(routerVars["id"])
	utils.TraceError(span, err)

	spending, err := handler.spending_repo.GetSpendingByUUId(context, spendingUUId)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if spending == nil {
		http.Error(writer, "Record not found", http.StatusNotFound)
		return
	}

	err = handler.spending_repo.LoadSpendingCategory(context, spending)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mappers.MapSpending(spending)

	err = json.NewEncoder(writer).Encode(response)
	utils.TraceError(span, err)
}
