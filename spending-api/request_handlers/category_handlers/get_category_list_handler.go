package category_handlers

import (
	"net/http"
	"spending/mappers"
	"spending/repositories/category_repo"
	"spending/request_handlers"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

type getCategoryListHandler struct {
	category_repo category_repo.CategoryRepository
}

func NewGetCategoryListHandler(categoryRepo category_repo.CategoryRepository) request_handlers.RequestHandler {
	return &getCategoryListHandler{
		category_repo: categoryRepo,
	}
}

func (handler *getCategoryListHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "GetCategoryListHandler")
	defer span.End()

	categories, err := handler.category_repo.GetCategoryList(context, nil)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response := mappers.MapCategoryList(categories)

	err = utils.Encode(context, writer, http.StatusOK, response)
	utils.TraceError(span, err)
}
