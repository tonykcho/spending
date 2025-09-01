package spending_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
	"spending/repositories/spending_repo"
	"spending/utils"
	"time"

	"go.opentelemetry.io/otel"
)

type CreateSpendingRequest struct {
	Amount       float32   `json:"amount"`
	Remark       string    `json:"remark"`
	SpendingDate time.Time `json:"spending_date"`
	Category     string    `json:"category"`
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
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request
	validationErrors := validateRequest(command)
	if validationErrors != nil {
		http.Error(writer, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	// Create a SpendingRecord from the request
	record := models.SpendingRecord{
		Amount:       command.Amount,
		Remark:       command.Remark,
		SpendingDate: command.SpendingDate,
		Category:     command.Category,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert the record into the database
	id := spending_repo.InsertSpendingRecord(record)

	spending := spending_repo.GetSpendingById(context, id)
	response := mappers.MapSpending(*spending)

	// Return 201 created response
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(response)
	utils.CheckError(err)
}

func validateRequest(request CreateSpendingRequest) error {
	if request.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if request.SpendingDate.IsZero() {
		return fmt.Errorf("spending date cannot be empty")
	}
	if request.Category == "" {
		return fmt.Errorf("category cannot be empty")
	}
	return nil
}
