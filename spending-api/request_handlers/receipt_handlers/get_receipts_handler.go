package receipt_handlers

import (
	"net/http"
	"spending/mappers"
	"spending/repositories/receipt_repo"
	"spending/request_handlers"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

type getReceiptsHandler struct {
	receipt_repo receipt_repo.ReceiptRepository
}

func NewGetReceiptsHandler(receiptRepo receipt_repo.ReceiptRepository) request_handlers.RequestHandler {
	return &getReceiptsHandler{
		receipt_repo: receiptRepo,
	}
}

func (handler *getReceiptsHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "GetReceiptsHandler")
	defer span.End()

	receipts, err := handler.receipt_repo.GetReceipts(ctx, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handler.receipt_repo.LoadReceiptsItems(ctx, nil, receipts)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mappers.MapReceipts(receipts)
	err = utils.Encode(ctx, writer, http.StatusOK, response)
	utils.TraceError(span, err)
}
