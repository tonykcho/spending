package spending_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
	"spending/repositories/category_repo"
	"spending/repositories/spending_repo"
	"spending/utils"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type CreateSpendingRequest struct {
	Amount       float32   `json:"amount"`
	Remark       string    `json:"remark"`
	SpendingDate time.Time `json:"spending_date"`
	CategoryId   uuid.UUID `json:"category_id"`
}

func CreateSpendingRequestHandler(writer http.ResponseWriter, request *http.Request) {
	// Trace the request
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "CreateSpendingRequestHandler")
	defer span.End()

	// Parse the request body into CreateSpendingRequest struct
	var command CreateSpendingRequest
	err := json.NewDecoder(request.Body).Decode(&command)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request
	validationErrors := validateRequest(command)
	if validationErrors != nil {
		utils.TraceError(span, validationErrors)
		http.Error(writer, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	category, err := category_repo.GetCategoryByUUId(context, command.CategoryId)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if category == nil {
		utils.TraceError(span, fmt.Errorf("category not found"))
		http.Error(writer, "category not found", http.StatusBadRequest)
		return
	}

	// Create a SpendingRecord from the request
	newSpending := models.SpendingRecord{
		Amount:       command.Amount,
		Remark:       command.Remark,
		SpendingDate: command.SpendingDate,
		CategoryId:   category.Id,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert the record into the database
	id := spending_repo.InsertSpendingRecord(context, newSpending)

	spending, err := spending_repo.GetSpendingById(context, id)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	response := mappers.MapSpending(spending)

	// Return 201 created response
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(response)
	utils.TraceError(span, err)
}

func validateRequest(request CreateSpendingRequest) error {
	if request.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if request.SpendingDate.IsZero() {
		return fmt.Errorf("spending date cannot be empty")
	}
	if request.CategoryId == uuid.Nil {
		return fmt.Errorf("category_id cannot be empty")
	}
	return nil
}
