package receipt_handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
	"spending/repositories"
	"spending/repositories/receipt_item_repo"
	"spending/repositories/receipt_repo"
	"spending/utils"
	"time"

	"go.opentelemetry.io/otel"
)

type createReceiptHandler struct {
	receipt_repo      receipt_repo.ReceiptRepository
	receipt_item_repo receipt_item_repo.ReceiptItemRepository
	unit_of_work      repositories.UnitOfWork
}

func NewCreateReceiptHandler(receiptRepo receipt_repo.ReceiptRepository, receiptItemRepo receipt_item_repo.ReceiptItemRepository, unitOfWork repositories.UnitOfWork) *createReceiptHandler {
	return &createReceiptHandler{
		receipt_repo:      receiptRepo,
		receipt_item_repo: receiptItemRepo,
		unit_of_work:      unitOfWork,
	}
}

type CreateReceiptRequest struct {
	StoreName   string                     `json:"storeName"`
	Date        time.Time                  `json:"date"`
	TotalAmount float64                    `json:"totalAmount"`
	Items       []CreateReceiptItemRequest `json:"items"`
}

type CreateReceiptItemRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (request CreateReceiptRequest) Valid(context context.Context) error {
	if request.StoreName == "" {
		return fmt.Errorf("store name cannot be empty")
	}
	if request.Date.IsZero() {
		return fmt.Errorf("date cannot be empty")
	}
	if request.TotalAmount < 0 {
		return fmt.Errorf("total amount cannot be negative")
	}
	for _, item := range request.Items {
		if err := item.Valid(context); err != nil {
			return fmt.Errorf("invalid item: %w", err)
		}
	}
	return nil
}

func (request CreateReceiptItemRequest) Valid(context context.Context) error {
	if request.Name == "" {
		return fmt.Errorf("item name cannot be empty")
	}
	if request.Price < 0 {
		return fmt.Errorf("item price cannot be negative")
	}
	return nil
}

func (handler *createReceiptHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "CreateReceiptHandler")
	defer span.End()

	command, err := utils.DecodeValid[CreateReceiptRequest](ctx, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var receipt *models.Receipt

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		var txErr error
		receipt = models.NewReceipt(command.StoreName, command.TotalAmount, command.Date)
		receipt, txErr = handler.receipt_repo.InsertReceipt(ctx, tx, receipt)
		if txErr != nil {
			return txErr
		}

		for _, item := range command.Items {
			receiptItem := models.NewReceiptItem(receipt.Id, item.Name, item.Price)
			receiptItem, txErr := handler.receipt_item_repo.InsertReceiptItem(ctx, tx, receiptItem)
			if txErr != nil {
				return txErr
			}

			receipt.Items = append(receipt.Items, receiptItem)
		}

		return nil
	})

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), utils.MapErrorToStatusCode(err))
		return
	}

	response := mappers.MapReceipt(receipt)
	writer.Header().Set("Location", fmt.Sprintf("/receipts/%s", receipt.UUId))
	err = utils.Encode(ctx, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
