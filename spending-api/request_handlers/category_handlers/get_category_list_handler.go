package category_handlers

import (
	"encoding/json"
	"net/http"
	"spending/mappers"
	"spending/repositories/category_repo"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

func GetCategoryListHandler(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "GetCategoryListHandler")
	defer span.End()

	categories := category_repo.GetCategoryList(context)
	response := mappers.MapCategoryList(categories)

	err := json.NewEncoder(writer).Encode(response)
	utils.TraceError(span, err)
}
