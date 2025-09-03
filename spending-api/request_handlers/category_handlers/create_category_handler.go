package category_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
	"spending/repositories/category_repo"
	"spending/utils"
	"time"

	"go.opentelemetry.io/otel"
)

type CreateCategoryHandler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
}

type createCategoryHandler struct {
	category_repo category_repo.CategoryRepository
}

func NewCreateCategoryHandler(categoryRepo category_repo.CategoryRepository) CreateCategoryHandler {
	return &createCategoryHandler{
		category_repo: categoryRepo,
	}
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (handler *createCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "CreateCategoryHandler")
	defer span.End()

	// Parse the request body into CreateCategoryRequest struct
	var command CreateCategoryRequest
	err := json.NewDecoder(request.Body).Decode(&command)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request
	validationErrors := validateCreateCategoryRequest(command)
	if validationErrors != nil {
		utils.TraceError(span, validationErrors)
		http.Error(writer, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	existingCategory, err := handler.category_repo.GetCategoryByName(context, command.Name)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if existingCategory != nil {
		utils.TraceError(span, fmt.Errorf("category already exists"))
		http.Error(writer, "Category already exists", http.StatusConflict)
		return
	}

	newCategory := models.Category{
		Name:      command.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := handler.category_repo.InsertCategory(context, newCategory)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := handler.category_repo.GetCategoryById(context, id)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response := mappers.MapCategory(category)

	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(response)
	utils.TraceError(span, err)
}

func validateCreateCategoryRequest(request CreateCategoryRequest) error {
	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
