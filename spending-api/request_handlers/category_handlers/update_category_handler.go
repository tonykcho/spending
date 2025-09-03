package category_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spending/repositories/category_repo"
	"spending/utils"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type UpdateCategoryHandler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
}

type updateCategoryHandler struct {
	category_repo category_repo.CategoryRepository
}

func NewUpdateCategoryHandler(categoryRepo category_repo.CategoryRepository) UpdateCategoryHandler {
	return &updateCategoryHandler{
		category_repo: categoryRepo,
	}
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
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

	var command UpdateCategoryRequest
	err = json.NewDecoder(request.Body).Decode(&command)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	validationErrors := validateUpdateCategoryRequest(command)
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

	if existingCategory != nil && existingCategory.UUId != categoryUUId {
		utils.TraceError(span, fmt.Errorf("category already exists"))
		http.Error(writer, "Category already exists", http.StatusConflict)
		return
	}

	category, err := handler.category_repo.GetCategoryByUUId(context, categoryUUId)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if category == nil {
		utils.TraceError(span, fmt.Errorf("category not found"))
		http.Error(writer, "Category not found", http.StatusNotFound)
		return
	}

	category.Name = command.Name
	category.UpdatedAt = time.Now()
	handler.category_repo.UpdateCategory(context, category)

	writer.WriteHeader(http.StatusNoContent)
}

func validateUpdateCategoryRequest(request UpdateCategoryRequest) error {
	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
