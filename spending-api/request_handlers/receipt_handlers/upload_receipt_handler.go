package receipt_handlers

import (
	"fmt"
	"net/http"
	"spending/dto"
	"spending/external_clients"
	"spending/utils"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type uploadReceiptHandler struct {
	paddle_ocr_client external_clients.PaddleOcrClient
	ollama_client     external_clients.OllamaClient
}

func NewUploadReceiptHandler(paddleOcrClient external_clients.PaddleOcrClient, ollamaClient external_clients.OllamaClient) *uploadReceiptHandler {
	return &uploadReceiptHandler{
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

	log.Info().Msg("Received file for OCR processing")

	ocrResult, err := handler.paddle_ocr_client.SendPaddleOcrRequest(ctx, file)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info().Msg("OCR processing completed, sending text to Llama3")

	ollamaResult, err := handler.ollama_client.GetJsonFromReceiptTextFromLLama3(ctx, ocrResult)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info().Msg("Successfully processed receipt and obtained structured data")

	result, err := processOllamaResult(ollamaResult)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, "Failed to process Ollama result: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.Encode(ctx, writer, http.StatusOK, result)
	utils.TraceError(span, err)
}

func processOllamaResult(result string) (*dto.ReceiptOcrDto, error) {
	// The receipt result is in format: store, item1:price1, item2:price2
	parts := strings.Split(result, "|")
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid result format")
	}

	store := parts[0]
	dateString := parts[1]
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		log.Info().Msgf("Error parsing date %s: %v", dateString, err)
		date = time.Now()
	}
	items := make([]dto.ReceiptItemOcrDto, 0)

	for _, item := range parts[2:] {
		itemParts := strings.Split(item, ":")
		if len(itemParts) != 2 {
			return nil, fmt.Errorf("invalid item format: %s", item)
		}

		itemName := itemParts[0]
		itemPrice, err := strconv.ParseFloat(sanitizePriceText(itemParts[1]), 64)
		if err != nil {
			log.Info().Msgf("Error parsing price for item %s: %v", itemName, err)
			continue
		}

		items = append(items, dto.ReceiptItemOcrDto{
			Name:  itemName,
			Price: itemPrice,
		})
	}

	// Process the extracted store and items as needed
	log.Info().Msgf("Extracted store: %s", store)
	for _, item := range items {
		log.Info().Msgf("Extracted item: %s, price: %.2f", item.Name, item.Price)
	}

	resultDto := &dto.ReceiptOcrDto{
		StoreName: store,
		Date:      date,
		Items:     items,
	}

	return resultDto, nil
}

func sanitizePriceText(priceText string) string {
	// Ai processed text may append text after the price, e.g., "23.5\ntotal"
	// This function extracts the numeric part only
	parts := strings.Fields(priceText)
	if len(parts) == 0 {
		return "0"
	}

	return parts[0]
}
