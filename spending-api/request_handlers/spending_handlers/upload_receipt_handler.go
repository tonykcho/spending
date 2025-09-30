package spending_handlers

import (
	"net/http"
	"spending/external_clients"
	"spending/repositories/spending_repo"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

type uploadReceiptHandler struct {
	spending_repo     spending_repo.SpendingRepository
	paddle_ocr_client external_clients.PaddleOcrClient
	ollama_client     external_clients.OllamaClient
}

func NewUploadReceiptHandler(spendingRepo spending_repo.SpendingRepository, paddleOcrClient external_clients.PaddleOcrClient, ollamaClient external_clients.OllamaClient) *uploadReceiptHandler {
	return &uploadReceiptHandler{
		spending_repo:     spendingRepo,
		paddle_ocr_client: paddleOcrClient,
		ollama_client:     ollamaClient,
	}
}

func (handler *uploadReceiptHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "UploadReceiptHandler")
	defer span.End()

	contentType := request.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(writer, "Content-Type header is missing", http.StatusBadRequest)
		return
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, "Failed to get file from form data: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ocrResult, err := handler.paddle_ocr_client.SendPaddleOcrRequest(ctx, file)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = utils.Encode(ctx, writer, http.StatusOK, ocrResult)
	utils.TraceError(span, err)
}
