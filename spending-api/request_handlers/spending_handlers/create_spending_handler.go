package spending_handlers

import (
	"context"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
	"spending/repositories"
	"spending/repositories/category_repo"
	"spending/repositories/spending_repo"
	"spending/request_handlers"
	"spending/utils"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type createSpendingHandler struct {
	spending_repo spending_repo.SpendingRepository
	category_repo category_repo.CategoryRepository
	unit_of_work  repositories.UnitOfWork
}

func NewCreateSpendingHandler(spendingRepo spending_repo.SpendingRepository, categoryRepo category_repo.CategoryRepository, unitOfWork repositories.UnitOfWork) request_handlers.RequestHandler {
	return &createSpendingHandler{
		spending_repo: spendingRepo,
		category_repo: categoryRepo,
		unit_of_work:  unitOfWork,
	}
}

type CreateSpendingRequest struct {
	Amount       float32   `json:"amount"`
	Remark       string    `json:"remark"`
	SpendingDate time.Time `json:"spending_date"`
	CategoryId   uuid.UUID `json:"category_id"`
}

func (request CreateSpendingRequest) Valid(context context.Context) error {
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

func (handler *createSpendingHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	// Trace the request
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "CreateSpendingHandler")
	defer span.End()

	// Parse the request body into CreateSpendingRequest struct
	command, err := utils.DecodeValid[CreateSpendingRequest](context, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := handler.unit_of_work.BeginTx()
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer handler.unit_of_work.CommitOrRollback(tx, err)

	category, err := handler.category_repo.GetCategoryByUUId(context, tx, command.CategoryId)
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
	id, err := handler.spending_repo.InsertSpendingRecord(context, tx, newSpending)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	spending, err := handler.spending_repo.GetSpendingById(context, tx, id)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	response := mappers.MapSpending(spending)

	// Return 201 created response
	err = utils.Encode(context, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
