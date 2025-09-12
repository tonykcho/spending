package spending_handlers

import (
	"context"
	"database/sql"
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
	SpendingDate time.Time `json:"spendingDate"`
	CategoryId   uuid.UUID `json:"categoryId"`
}

func (request CreateSpendingRequest) Valid(context context.Context) error {
	if request.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if request.SpendingDate.IsZero() {
		return fmt.Errorf("spending date cannot be empty")
	}
	if request.CategoryId == uuid.Nil {
		return fmt.Errorf("categoryId cannot be empty")
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

	var spending *models.SpendingRecord

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		category, txErr := handler.category_repo.GetCategoryByUUId(context, tx, command.CategoryId)
		if txErr != nil {
			return txErr
		}

		if category == nil {
			txErr := fmt.Errorf("category not found")
			return txErr
		}

		// Create a SpendingRecord from the request
		newSpending := models.NewSpendingRecord(command.Amount, command.Remark, command.SpendingDate, category.Id)

		// Insert the record into the database
		spending, txErr = handler.spending_repo.InsertSpendingRecord(context, tx, newSpending)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), utils.MapErrorToStatusCode(err))
		return
	}

	response := mappers.MapSpending(spending)
	writer.Header().Set("Location", fmt.Sprintf("/spending/%s", spending.UUId))
	err = utils.Encode(context, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
