package category_handlers

import (
	"context"
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
	context, span := tracer.Start(request.Context(), "CreateCategoryHandler")
	defer span.End()

	// Parse the request body into CreateCategoryRequest struct
	command, err := utils.DecodeValid[CreateCategoryRequest](context, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Start a new transaction
	tx, err := handler.unit_of_work.BeginTx()
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	existingCategory, err := handler.category_repo.GetCategoryByName(context, tx, command.Name)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if existingCategory != nil {
		err = fmt.Errorf("category already exists")
		utils.TraceError(span, fmt.Errorf("category already exists"))
		http.Error(writer, "Category already exists", http.StatusConflict)
		return
	}

	newCategory := models.Category{
		Name:      command.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := handler.category_repo.InsertCategory(context, tx, newCategory)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	category, err := handler.category_repo.GetCategoryById(context, tx, id)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response := mappers.MapCategory(category)
	err = utils.Encode(context, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
