package spending_handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSpendingRequestHandler(t *testing.T) {
	// Prepare a valid JSON request body
	body := []byte(`{
        "amount": 100,
        "remark": "Test",
        "spending_date": "2023-10-01T00:00:00Z",
        "category": "Food"
    }`)

	request := httptest.NewRequest(http.MethodPost, "/spending", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	CreateSpendingRequestHandler(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", recorder.Code)
	}

	// Optionally, check the response body for expected fields
	t.Log(recorder.Body.String())
}
