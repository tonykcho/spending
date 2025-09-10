package category_handlers

import (
	"net/http"
	"spending/mappers"
	"spending/repositories/category_repo"
	"spending/request_handlers"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type getCategoryHandler struct {
	category_repo category_repo.CategoryRepository
}

func NewGetCategoryHandler(categoryRepo category_repo.CategoryRepository) request_handlers.RequestHandler {
	return &getCategoryHandler{
		category_repo: categoryRepo,
	}
}

func (handler *getCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "GetCategoryHandler")
	defer span.End()

	routerVars := mux.Vars(request)

	categoryUUId, err := uuid.Parse(routerVars["id"])
	utils.TraceError(span, err)

	category, err := handler.category_repo.GetCategoryByUUId(ctx, nil, categoryUUId)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if category == nil {
		http.Error(writer, "Record not found", http.StatusNotFound)
		return
	}

	err = handler.category_repo.LoadStoresForCategory(ctx, nil, category)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mappers.MapCategory(category)

	err = utils.Encode(ctx, writer, http.StatusOK, response)
	utils.TraceError(span, err)
}
