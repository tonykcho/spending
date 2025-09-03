package category_handlers

import (
	"encoding/json"
	"net/http"
	"spending/mappers"
	"spending/repositories/category_repo"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

type GetCategoryListHandler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
}

type getCategoryListHandler struct {
	category_repo category_repo.CategoryRepository
}

func NewGetCategoryListHandler(categoryRepo category_repo.CategoryRepository) GetCategoryListHandler {
	return &getCategoryListHandler{
		category_repo: categoryRepo,
	}
}

func (handler *getCategoryListHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "GetCategoryListHandler")
	defer span.End()

	categories, err := handler.category_repo.GetCategoryList(context)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response := mappers.MapCategoryList(categories)

	err = json.NewEncoder(writer).Encode(response)
	utils.TraceError(span, err)
}
