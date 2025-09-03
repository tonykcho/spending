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

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func CreateCategoryRequestHandler(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "CreateCategoryRequestHandler")
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

	existingCategory := category_repo.GetCategoryByName(context, command.Name)
	if existingCategory != nil {
		utils.TraceError(span, fmt.Errorf("category already exists"))
		http.Error(writer, "Category already exists", http.StatusConflict)
		return
	}

	category := models.Category{
		Name:      command.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id := category_repo.InsertCategory(context, category)

	category = *category_repo.GetCategoryById(context, id)
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
