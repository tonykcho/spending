package category_handlers

import (
	"database/sql"
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

	status := http.StatusInternalServerError

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		txErr := handler.category_repo.DeleteCategory(context, nil, categoryUUId)
		if txErr != nil {
			status = http.StatusInternalServerError
			return txErr
		}

		return nil
	})

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), status)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
