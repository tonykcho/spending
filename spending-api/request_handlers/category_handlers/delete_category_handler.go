package category_handlers

import (
	"net/http"
	"spending/repositories/category_repo"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type DeleteCategoryHandler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
}

type deleteCategoryHandler struct {
	category_repo category_repo.CategoryRepository
}

func NewDeleteCategoryHandler(categoryRepo category_repo.CategoryRepository) DeleteCategoryHandler {
	return &deleteCategoryHandler{
		category_repo: categoryRepo,
	}
}

func (handler *deleteCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "DeleteCategoryHandler")
	defer span.End()

	routerVars := mux.Vars(request)
	categoryUUId, err := uuid.Parse(routerVars["id"])
	utils.TraceError(span, err)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handler.category_repo.DeleteCategory(context, categoryUUId)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
