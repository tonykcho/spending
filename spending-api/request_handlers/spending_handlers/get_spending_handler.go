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

func GetSpendingRequestHandler(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "DB:GetSpendingByUUId")
	defer span.End()

	routerVars := mux.Vars(request)

	spendingUUId, err := uuid.Parse(routerVars["id"])
	utils.CheckError(err)

	spending := spending_repo.GetSpendingByUUId(context, spendingUUId)

	if spending == nil {
		http.Error(writer, "Record not found", http.StatusNotFound)
		return
	}

	response := mappers.MapSpending(*spending)

	err = json.NewEncoder(writer).Encode(response)
	utils.CheckError(err)
}
