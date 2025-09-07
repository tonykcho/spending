package category_handlers

import (
	"net/http"
	"spending/repositories"
	"spending/repositories/category_repo"
	"spending/request_handlers"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type deleteCategoryHandler struct {
	category_repo category_repo.CategoryRepository
	unit_of_work  repositories.UnitOfWork
}

func NewDeleteCategoryHandler(categoryRepo category_repo.CategoryRepository, unitOfWork repositories.UnitOfWork) request_handlers.RequestHandler {
	return &deleteCategoryHandler{
		category_repo: categoryRepo,
		unit_of_work:  unitOfWork,
	}
}

func (handler *deleteCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "DeleteCategoryHandler")
	defer span.End()

	routerVars := mux.Vars(request)
	categoryUUId, err := uuid.Parse(routerVars["id"])
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Start a new transaction
	tx, err := handler.unit_of_work.BeginTx()
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer handler.unit_of_work.CommitOrRollback(tx, err)

	err = handler.category_repo.DeleteCategory(context, tx, categoryUUId)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
