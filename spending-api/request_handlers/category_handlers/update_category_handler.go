package category_handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"spending/repositories"
	"spending/repositories/category_repo"
	"spending/request_handlers"
	"spending/utils"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type updateCategoryHandler struct {
	category_repo category_repo.CategoryRepository
	unit_of_work  repositories.UnitOfWork
}

func NewUpdateCategoryHandler(categoryRepo category_repo.CategoryRepository, unitOfWork repositories.UnitOfWork) request_handlers.RequestHandler {
	return &updateCategoryHandler{
		category_repo: categoryRepo,
		unit_of_work:  unitOfWork,
	}
}

type UpdateCategoryRequest struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (request UpdateCategoryRequest) Valid(context context.Context) error {
	if request.Id == uuid.Nil {
		return fmt.Errorf("id cannot be empty")
	}

	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

func (handler *updateCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "UpdateCategoryHandler")
	defer span.End()

	routerVars := mux.Vars(request)
	categoryUUId, err := uuid.Parse(routerVars["id"])
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	command, err := utils.DecodeValid[UpdateCategoryRequest](context, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if command.Id != categoryUUId {
		err = fmt.Errorf("id in path and body do not match")
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	status := http.StatusInternalServerError

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		existingCategory, txErr := handler.category_repo.GetCategoryByName(context, tx, command.Name)
		if txErr != nil {
			status = http.StatusInternalServerError
			return txErr
		}

		if existingCategory != nil && existingCategory.UUId != categoryUUId {
			txErr = fmt.Errorf("category with name %s already exists", command.Name)
			status = http.StatusConflict
			return txErr
		}

		category, txErr := handler.category_repo.GetCategoryByUUId(context, tx, categoryUUId)
		if txErr != nil {
			status = http.StatusInternalServerError
			return txErr
		}

		if category == nil {
			txErr = fmt.Errorf("category not found")
			status = http.StatusNotFound
			return txErr
		}

		category.Name = command.Name
		category.UpdatedAt = time.Now()
		txErr = handler.category_repo.UpdateCategory(context, tx, category)
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
