package external_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

type OllamaClient interface {
	GetJsonFromReceiptTextFromLLama3(ctx context.Context, texts []string) (string, error)
}

type ollamaClient struct {
}

func NewOllamaClient() OllamaClient {
	return &ollamaClient{}
}

// GetJsonFromReceiptTextFromLLama3 sends the prompt to the local Ollama Llama3 model and returns the response as a string
func (c *ollamaClient) GetJsonFromReceiptTextFromLLama3(ctx context.Context, texts []string) (string, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "GetJsonFromReceiptTextFromLLama3")
	defer span.End()

	ollamaHost := utils.GetOllamaHost()
	url := fmt.Sprintf("%s/api/chat", ollamaHost)

	// Construct the chat request payload
	messages := []map[string]string{
		{"role": "system", "content": "You are a helpful assistant that extracts structured data from receipts."},
		{"role": "system", "content": "Respond only with JSON, no explanations."},
		{"role": "system", "content": "Determine if the texts are receipts. If not, respond with {\"error\": \"not a receipt\"}."},
		{"role": "system", "content": "If the texts are receipts, extract the following fields: store (string), date (string in YYYY-MM-DD format), total (number), items (array of objects with name (string), price (number), discount (number, optional))."},
		{"role": "system", "content": "If any field is missing or cannot be determined, use null for that field."},
		{"role": "system", "content": "Here is an example response for a valid receipt: {\"store\": \"Store Name\", \"date\": \"2023-10-01\", \"total\": 23.45, \"items\": [{\"name\": \"Item 1\", \"price\": 10.00}, {\"name\": \"Item 2\", \"price\": 15.00, \"discount\": 1.55}]}."},
		{"role": "system", "content": "Preserve all the chinese characters in the receipt text."},
		{"role": "user", "content": fmt.Sprintf("Extract the following receipt texts into JSON format with fields: store, date, total, items (with name, price, discount). Here are the texts: %v", texts)},
	}

	payload := map[string]interface{}{
		"model":    "llama3",
		"messages": messages,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		utils.TraceError(span, err)
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		utils.TraceError(span, err)
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.TraceError(span, err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 response: %d", response.StatusCode)
		utils.TraceError(span, err)
		return "", err
	}

	return "", nil
}
