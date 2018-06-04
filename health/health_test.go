package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Health)

	handler.ServeHTTP(rr, req)

	// Checking Status Code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned incorrect code: Received %v, expected %v", status, http.StatusOK)
	}

	// Checking Response Value
	expected := `"Ok"`
	if rr.Body.String() != expected {
		t.Errorf("Received unexpected response: Received %v, expected %v", rr.Body.String(), expected)
	}
}
