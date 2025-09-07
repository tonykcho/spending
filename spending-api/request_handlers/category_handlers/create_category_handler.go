package category_handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
	"spending/repositories"
	"spending/repositories/category_repo"
	"spending/request_handlers"
	"spending/utils"
	"time"

	"go.opentelemetry.io/otel"
)

type createCategoryHandler struct {
	category_repo category_repo.CategoryRepository
	unit_of_work  repositories.UnitOfWork
}

func NewCreateCategoryHandler(categoryRepo category_repo.CategoryRepository, unitOfWork repositories.UnitOfWork) request_handlers.RequestHandler {
	return &createCategoryHandler{
		category_repo: categoryRepo,
		unit_of_work:  unitOfWork,
	}
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (request CreateCategoryRequest) Valid(context context.Context) error {
	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

func (handler *createCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "CreateCategoryHandler")
	defer span.End()

	// Parse the request body into CreateCategoryRequest struct
	command, err := utils.DecodeValid[CreateCategoryRequest](ctx, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var category *models.Category
	var status int = http.StatusInternalServerError

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		existingCategory, txErr := handler.category_repo.GetCategoryByName(ctx, tx, command.Name)
		if txErr != nil {
			status = http.StatusInternalServerError
			return txErr
		}

		if existingCategory != nil {
			txErr = fmt.Errorf("category already exists")
			status = http.StatusConflict
			return txErr
		}

		newCategory := models.Category{
			Name:      command.Name,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		category, txErr = handler.category_repo.InsertCategory(ctx, tx, newCategory)
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

	response := mappers.MapCategory(category)
	writer.Header().Set("Location", fmt.Sprintf("/categories/%s", category.UUId))
	err = utils.Encode(ctx, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
