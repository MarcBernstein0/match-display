package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckController(t *testing.T) {
	// Create a request to pass our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter
	req, err := http.NewRequest(http.MethodGet, "/match-display/v1/health-check", nil)
	if err != nil {
		t.Fatalf("error in new request\n%v", err)
	}

	// Create a ResponseRecorder (which satisfies httpResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our requst and response recoder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: Expected=%d got=%d", http.StatusOK, status)
	}

	// Check the response body is what we expect
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: Expected=%s, got=%s",
			expected, rr.Body.String())
	}
}

func TestTournamentController(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/match-display/v1/tournaments", nil)
	if err != nil {
		t.Fatalf("error in new request\n%v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTournamentData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: Expected=%d got=%d", http.StatusOK, status)
	}

}
