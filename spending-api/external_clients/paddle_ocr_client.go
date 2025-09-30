package external_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

// Send http post request with jpeg/png content type to paddle ocr server, get back the ocr string array result
func SendPaddleOcrRequest(ctx context.Context, file multipart.File) ([]string, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(ctx, "SendPaddleOcrRequest")
	defer span.End()

	paddleOcrHost := utils.GetPaddleOcrHost()
	url := fmt.Sprintf("%s/upload_image/", paddleOcrHost)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	contentType := http.DetectContentType(fileBytes[:min(len(fileBytes), 512)])
	if contentType != "image/jpeg" && contentType != "image/png" {
		err = fmt.Errorf("unsupported content type: %s", contentType)
		utils.TraceError(span, err)
		return nil, err
	}

	var bodyBuf bytes.Buffer
	writer := multipart.NewWriter(&bodyBuf)
	part, err := writer.CreateFormFile("file", "receipt"+extForContentType(contentType))
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}
	if _, err = part.Write(fileBytes); err != nil {
		utils.TraceError(span, err)
		return nil, err
	}
	_ = writer.Close()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &bodyBuf)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("paddle ocr server returned status code: %d", response.StatusCode)
		utils.TraceError(span, err)
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	var ocrResult []string
	err = json.Unmarshal(body, &ocrResult)
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	return ocrResult, nil
}

func extForContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ".jpg"
	}
}
