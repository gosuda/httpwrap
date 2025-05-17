package httperror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestNewRFC7807Error(t *testing.T) {
	tests := []struct {
		name   string
		status int
		title  string
		detail string
	}{
		{
			name:   "400 Bad Request",
			status: 400,
			title:  "Bad Request",
			detail: "Invalid input parameters",
		},
		{
			name:   "404 Not Found",
			status: 404,
			title:  "Not Found",
			detail: "Resource could not be located",
		},
		{
			name:   "500 Internal Server Error",
			status: 500,
			title:  "Internal Server Error",
			detail: "An unexpected error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC7807Error(tt.status, tt.title, tt.detail)

			if err.Status != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, err.Status)
			}

			if err.Title != tt.title {
				t.Errorf("Expected title %s, got %s", tt.title, err.Title)
			}

			if err.Detail != tt.detail {
				t.Errorf("Expected detail %s, got %s", tt.detail, err.Detail)
			}

			if err.Type != "about:blank" {
				t.Errorf("Expected default type about:blank, got %s", err.Type)
			}
		})
	}
}

func TestRFC7807Error_Error(t *testing.T) {
	tests := []struct {
		name   string
		status int
		title  string
		detail string
		want   string
	}{
		{
			name:   "400 Bad Request",
			status: 400,
			title:  "Bad Request",
			detail: "Invalid input parameters",
			want:   "400: Bad Request - Invalid input parameters",
		},
		{
			name:   "404 Not Found",
			status: 404,
			title:  "Not Found",
			detail: "Resource could not be located",
			want:   "404: Not Found - Resource could not be located",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC7807Error(tt.status, tt.title, tt.detail)

			if got := err.Error(); got != tt.want {
				t.Errorf("RFC7807Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRFC7807Error_WithType(t *testing.T) {
	err := NewRFC7807Error(400, "Bad Request", "Invalid parameters")
	typeURI := "https://example.com/errors/validation"

	err = err.WithType(typeURI)

	if err.Type != typeURI {
		t.Errorf("Expected Type %s, got %s", typeURI, err.Type)
	}
}

func TestRFC7807Error_WithInstance(t *testing.T) {
	err := NewRFC7807Error(404, "Not Found", "Resource not found")
	instance := "/api/resources/123"

	err = err.WithInstance(instance)

	if err.Instance != instance {
		t.Errorf("Expected Instance %s, got %s", instance, err.Instance)
	}
}

func TestRFC7807Error_WithExtension(t *testing.T) {
	err := NewRFC7807Error(403, "Forbidden", "Insufficient permissions")

	// Add extensions
	err = err.WithExtension("resource_id", "123")
	err = err.WithExtension("required_role", "admin")

	if val, ok := err.Extensions["resource_id"]; !ok || val != "123" {
		t.Errorf("Expected extension resource_id=123, got %v", val)
	}

	if val, ok := err.Extensions["required_role"]; !ok || val != "admin" {
		t.Errorf("Expected extension required_role=admin, got %v", val)
	}
}

func TestRFC7807Error_MarshalJSON(t *testing.T) {
	err := NewRFC7807Error(403, "Forbidden", "Insufficient permissions")
	err = err.WithType("https://example.com/errors/forbidden")
	err = err.WithInstance("/api/resources/123")
	err = err.WithExtension("resource_id", "123")
	err = err.WithExtension("required_role", "admin")

	bytes, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal RFC7807Error: %v", marshalErr)
	}

	// Unmarshal to check content
	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(bytes, &result); unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal JSON result: %v", unmarshalErr)
	}

	// Check standard fields
	if result["type"] != "https://example.com/errors/forbidden" {
		t.Errorf("Expected type field, got %v", result["type"])
	}

	if result["title"] != "Forbidden" {
		t.Errorf("Expected title field, got %v", result["title"])
	}

	if result["detail"] != "Insufficient permissions" {
		t.Errorf("Expected detail field, got %v", result["detail"])
	}

	if result["instance"] != "/api/resources/123" {
		t.Errorf("Expected instance field, got %v", result["instance"])
	}

	// Check extensions
	if result["resource_id"] != "123" {
		t.Errorf("Expected resource_id extension, got %v", result["resource_id"])
	}

	if result["required_role"] != "admin" {
		t.Errorf("Expected required_role extension, got %v", result["required_role"])
	}
}

func TestRFC7807Error_StatusCode(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"400 Bad Request", 400},
		{"404 Not Found", 404},
		{"500 Internal Server Error", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC7807Error(tt.status, "Test Title", "Test Detail")

			if got := err.StatusCode(); got != tt.status {
				t.Errorf("RFC7807Error.StatusCode() = %v, want %v", got, tt.status)
			}
		})
	}
}

func TestRFC7807Error_ErrorMessage(t *testing.T) {
	tests := []struct {
		name   string
		detail string
	}{
		{"Simple message", "Simple error message"},
		{"Empty message", ""},
		{"Message with special characters", "Error: $#@! special chars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC7807Error(400, "Test Title", tt.detail)

			if got := err.ErrorMessage(); got != tt.detail {
				t.Errorf("RFC7807Error.ErrorMessage() = %v, want %v", got, tt.detail)
			}
		})
	}
}

func TestRFC7807Error_ToHttpError(t *testing.T) {
	// Create a problem detail
	problem := NewRFC7807Error(403, "Forbidden", "Access denied")
	problem = problem.WithType("https://example.com/errors/forbidden")
	problem = problem.WithInstance("/api/resource/123")
	problem = problem.WithExtension("user_id", "abc123")

	// Convert to HttpError
	httpErr := problem.ToHttpError()

	// Check status code
	if httpErr.StatusCode() != 403 {
		t.Errorf("Expected status code 403, got %d", httpErr.StatusCode())
	}

	// Try to parse the message as JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(httpErr.ErrorMessage()), &result); err != nil {
		t.Fatalf("Failed to unmarshal error message as JSON: %v", err)
	}

	// Check JSON fields
	if result["type"] != "https://example.com/errors/forbidden" {
		t.Errorf("Expected type field, got %v", result["type"])
	}

	if result["title"] != "Forbidden" {
		t.Errorf("Expected title field, got %v", result["title"])
	}

	if result["status"].(float64) != 403 {
		t.Errorf("Expected status field to be 403, got %v", result["status"])
	}

	if result["detail"] != "Access denied" {
		t.Errorf("Expected detail field, got %v", result["detail"])
	}

	if result["instance"] != "/api/resource/123" {
		t.Errorf("Expected instance field, got %v", result["instance"])
	}

	if result["user_id"] != "abc123" {
		t.Errorf("Expected user_id extension, got %v", result["user_id"])
	}
}

func TestProblemHelperFunctions(t *testing.T) {
	tests := []struct {
		name      string
		fn        func(string, ...string) *RFC7807Error
		status    int
		detail    string
		wantTitle string
	}{
		{
			name:      "BadRequestProblem",
			fn:        BadRequestProblem,
			status:    http.StatusBadRequest,
			detail:    "Invalid input",
			wantTitle: "Bad Request",
		},
		{
			name:      "BadRequestProblem with custom title",
			fn:        BadRequestProblem,
			status:    http.StatusBadRequest,
			detail:    "Invalid input",
			wantTitle: "Validation Error",
		},
		{
			name:      "NotFoundProblem",
			fn:        NotFoundProblem,
			status:    http.StatusNotFound,
			detail:    "Resource not found",
			wantTitle: "Not Found",
		},
		{
			name:      "NotFoundProblem with custom title",
			fn:        NotFoundProblem,
			status:    http.StatusNotFound,
			detail:    "Resource not found",
			wantTitle: "Missing Resource",
		},
		{
			name:      "ForbiddenProblem",
			fn:        ForbiddenProblem,
			status:    http.StatusForbidden,
			detail:    "Access denied",
			wantTitle: "Forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err *RFC7807Error

			if tt.name == "BadRequestProblem with custom title" {
				err = tt.fn(tt.detail, "Validation Error")
			} else if tt.name == "NotFoundProblem with custom title" {
				err = tt.fn(tt.detail, "Missing Resource")
			} else {
				err = tt.fn(tt.detail)
			}

			if err.Status != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, err.Status)
			}

			if err.Detail != tt.detail {
				t.Errorf("Expected detail %s, got %s", tt.detail, err.Detail)
			}

			if err.Title != tt.wantTitle {
				t.Errorf("Expected title %s, got %s", tt.wantTitle, err.Title)
			}
		})
	}
}

func TestRFC7807CompleteFlow(t *testing.T) {
	// Create a problem with all fields
	problem := BadRequestProblem("Invalid user ID format")
	problem = problem.WithType("https://example.com/errors/validation")
	problem = problem.WithInstance("/api/users/abc")
	problem = problem.WithExtension("invalid_field", "user_id")
	problem = problem.WithExtension("provided_value", "abc")
	problem = problem.WithExtension("expected_format", "numeric")

	// Validate the Problem Detail fields
	if problem.Status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, problem.Status)
	}

	if problem.Title != "Bad Request" {
		t.Errorf("Expected title %s, got %s", "Bad Request", problem.Title)
	}

	if problem.Detail != "Invalid user ID format" {
		t.Errorf("Expected detail %s, got %s", "Invalid user ID format", problem.Detail)
	}

	if problem.Type != "https://example.com/errors/validation" {
		t.Errorf("Expected type %s, got %s", "https://example.com/errors/validation", problem.Type)
	}

	if problem.Instance != "/api/users/abc" {
		t.Errorf("Expected instance %s, got %s", "/api/users/abc", problem.Instance)
	}

	// Convert to HttpError
	httpErr := problem.ToHttpError()

	// Validate HttpError
	expectedErrorCode := 400
	if httpErr.StatusCode() != expectedErrorCode {
		t.Errorf("Expected HTTP error code %d, got %d", expectedErrorCode, httpErr.StatusCode())
	}

	// Parse message as JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(httpErr.ErrorMessage()), &result); err != nil {
		t.Fatalf("Failed to unmarshal error message as JSON: %v", err)
	}

	// Check all fields
	expectedFields := []struct {
		key  string
		want interface{}
	}{
		{"type", "https://example.com/errors/validation"},
		{"title", "Bad Request"},
		{"status", float64(400)},
		{"detail", "Invalid user ID format"},
		{"instance", "/api/users/abc"},
		{"invalid_field", "user_id"},
		{"provided_value", "abc"},
		{"expected_format", "numeric"},
	}

	for _, field := range expectedFields {
		if result[field.key] != field.want {
			t.Errorf("Expected %s = %v, got %v", field.key, field.want, result[field.key])
		}
	}

	// Test the Error() string
	expectedErrStr := "400: Bad Request - Invalid user ID format"
	if got := problem.Error(); got != expectedErrStr {
		t.Errorf("Expected Error() to be %q, got %q", expectedErrStr, got)
	}

	// Validate that the problem implements the error interface
	var stdErr error = problem
	if stdErr.Error() != expectedErrStr {
		t.Errorf("Problem does not properly implement error interface")
	}

	fmt.Println("RFC7807 complete flow test passed successfully")
}
