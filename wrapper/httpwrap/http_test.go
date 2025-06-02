package httpwrap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gosuda/httpwrap/httperror"
)

func TestMux_HandleWithContentType(t *testing.T) {
	tests := []struct {
		name                string
		handler             HandlerFunc
		expectedStatus      int
		expectedContentType string
		expectedBody        string
	}{
		{
			name: "HttpError with JSON content type",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return httperror.New(400, `{"error":"bad request"}`, "application/json")
			},
			expectedStatus:      400,
			expectedContentType: "application/json",
			expectedBody:        `{"error":"bad request"}`,
		},
		{
			name: "HttpError with plain text content type",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return httperror.New(404, "Not found", "text/plain")
			},
			expectedStatus:      404,
			expectedContentType: "text/plain",
			expectedBody:        "Not found",
		},
		{
			name: "HttpError without content type",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return httperror.New(500, "Internal server error")
			},
			expectedStatus:      500,
			expectedContentType: "text/plain; charset=utf-8", // Default from http.Error
			expectedBody:        "Internal server error",
		},
		{
			name: "RFC9457 error with application/problem+json",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				prob := httperror.BadRequestProblem9457("Invalid input data")
				return prob.ToHttpError()
			},
			expectedStatus:      400,
			expectedContentType: "application/problem+json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := NewMux(nil)
			mux.Handle("/test", tt.handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if tt.expectedContentType != "" && contentType != tt.expectedContentType {
				t.Errorf("Expected content type %s, got %s", tt.expectedContentType, contentType)
			}

			if tt.expectedBody != "" {
				actualBody := w.Body.String()
				expectedBody := tt.expectedBody

				// http.Error adds a newline, but direct write doesn't
				// Check if this error has ContentType (direct write) or not (http.Error)
				hasCustomContentType := false
				for _, tc := range []string{"application/json", "application/problem+json", "text/plain"} {
					if contentType == tc {
						hasCustomContentType = true
						break
					}
				}

				if !hasCustomContentType {
					expectedBody += "\n"
				}

				if actualBody != expectedBody {
					t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
				}
			}
		})
	}
}

func TestMux_HandleWithoutError(t *testing.T) {
	mux := NewMux(nil)
	mux.Handle("/test", func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
		return nil
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body OK, got %s", w.Body.String())
	}
}
