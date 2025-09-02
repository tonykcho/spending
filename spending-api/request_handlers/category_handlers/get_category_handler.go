package category_handlers

import (
	"encoding/json"
	"net/http"
	"spending/mappers"
	"spending/repositories/category_repo"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

func GetCategoryRequestHandler(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "GetCategoryRequestHandler")
	defer span.End()

	routerVars := mux.Vars(request)

	categoryUUId, err := uuid.Parse(routerVars["id"])
	utils.TraceError(span, err)

	category := category_repo.GetCategoryByUUId(context, categoryUUId)

	if category == nil {
		http.Error(writer, "Record not found", http.StatusNotFound)
		return
	}

	response := mappers.MapCategory(*category)

	err = json.NewEncoder(writer).Encode(response)
	utils.TraceError(span, err)
}
