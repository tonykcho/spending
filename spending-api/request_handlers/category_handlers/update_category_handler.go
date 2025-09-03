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

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

func UpdateCategoryRequestHandler(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "CreateCategoryRequestHandler")
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

	existingCategory := category_repo.GetCategoryByName(context, command.Name)
	if existingCategory != nil && existingCategory.UUId != categoryUUId {
		utils.TraceError(span, fmt.Errorf("category already exists"))
		http.Error(writer, "Category already exists", http.StatusConflict)
		return
	}

	category := category_repo.GetCategoryByUUId(context, categoryUUId)
	if category == nil {
		utils.TraceError(span, fmt.Errorf("category not found"))
		http.Error(writer, "Category not found", http.StatusNotFound)
		return
	}

	category.Name = command.Name
	category.UpdatedAt = time.Now()
	category_repo.UpdateCategory(context, category)

	writer.WriteHeader(http.StatusNoContent)
}

func validateUpdateCategoryRequest(request UpdateCategoryRequest) error {
	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
