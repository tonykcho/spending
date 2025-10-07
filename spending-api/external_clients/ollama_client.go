package external_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"spending/utils"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type OllamaResult struct {
	Model              string              `json:"model"`
	CreateAt           string              `json:"created_at"`
	Message            OllamaResultMessage `json:"message"`
	Done               bool                `json:"done"`
	DoneReason         string              `json:"done_reason"`
	TotalDuration      int64               `json:"total_duration"`
	LoadDuration       int64               `json:"load_duration"`
	PromptEvalCount    int64               `json:"prompt_eval_count"`
	PromptEvalDuration int64               `json:"prompt_eval_duration"`
	EvalCount          int64               `json:"eval_count"`
	EvalDuration       int64               `json:"eval_duration"`
}

type OllamaResultMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

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
		{
			"role":    "system",
			"content": "提取店舖名和所有收據貨品(格式：店舖名|日期|貨品1:價格1|貨品2:價格2)。不用解釋",
		},
		{
			"role":    "user",
			"content": fmt.Sprintf("Receipt Texts: %v", texts),
		},
	}

	// log.Info().Msgf("texts: %v", texts)
	log.Info().Msgf("Sending messages to Ollama: %v", messages)

	payload := map[string]interface{}{
		"model":    "llama3.1:8b",
		"options":  map[string]interface{}{"temperature": 0},
		"messages": messages,
		"stream":   false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		utils.TraceError(span, err)
		return "", err
	}

	log.Info().Msgf("Sending request to Ollama at %s", url)

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

	var ollamaResult OllamaResult
	err = json.NewDecoder(response.Body).Decode(&ollamaResult)
	if err != nil {
		utils.TraceError(span, err)
		return "", err
	}

	log.Info().Msgf("Ollama response received: %+v", ollamaResult)

	return ollamaResult.Message.Content, nil
}
