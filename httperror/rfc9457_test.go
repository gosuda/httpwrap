package httperror

import (
	"encoding/json"
	"net/http"
	"testing"
)

// Test basic RFC9457Error creation
func TestNewRFC9457Error(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		title    string
		detail   string
		wantType string
	}{
		{
			name:     "400 Bad Request",
			status:   http.StatusBadRequest,
			title:    "Bad Request",
			detail:   "The request is invalid",
			wantType: "about:blank",
		},
		{
			name:     "404 Not Found",
			status:   http.StatusNotFound,
			title:    "Not Found",
			detail:   "The requested resource was not found",
			wantType: "about:blank",
		},
		{
			name:     "500 Internal Server Error",
			status:   http.StatusInternalServerError,
			title:    "Internal Server Error",
			detail:   "An unexpected error occurred",
			wantType: "about:blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC9457Error(tt.status, tt.title, tt.detail)

			if err.Status != tt.status {
				t.Errorf("Status = %d, want %d", err.Status, tt.status)
			}

			if err.Title != tt.title {
				t.Errorf("Title = %s, want %s", err.Title, tt.title)
			}

			if err.Detail != tt.detail {
				t.Errorf("Detail = %s, want %s", err.Detail, tt.detail)
			}

			if err.Type != tt.wantType {
				t.Errorf("Type = %s, want %s", err.Type, tt.wantType)
			}
		})
	}
}

// Test RFC9457Error with custom type
func TestNewRFC9457ErrorWithType(t *testing.T) {
	status := http.StatusBadRequest
	typeURI := "https://example.com/problems/validation-error"
	title := "Validation Error"
	detail := "The input data is invalid"

	err := NewRFC9457ErrorWithType(status, typeURI, title, detail)

	if err.Status != status {
		t.Errorf("Status = %d, want %d", err.Status, status)
	}

	if err.Type != typeURI {
		t.Errorf("Type = %s, want %s", err.Type, typeURI)
	}

	if err.Title != title {
		t.Errorf("Title = %s, want %s", err.Title, title)
	}

	if err.Detail != detail {
		t.Errorf("Detail = %s, want %s", err.Detail, detail)
	}
}

// Test RFC9457Error.Error() method
func TestRFC9457Error_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *RFC9457Error
		want string
	}{
		{
			name: "400 Bad Request",
			err:  NewRFC9457Error(400, "Bad Request", "Invalid input"),
			want: "400: Bad Request - Invalid input",
		},
		{
			name: "404 Not Found",
			err:  NewRFC9457Error(404, "Not Found", "Resource not found"),
			want: "404: Not Found - Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("RFC9457Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test WithType method
func TestRFC9457Error_WithType(t *testing.T) {
	err := NewRFC9457Error(400, "Bad Request", "Invalid input")
	typeURI := "https://example.com/problems/bad-request"

	result := err.WithType(typeURI)

	if result.Type != typeURI {
		t.Errorf("Type = %s, want %s", result.Type, typeURI)
	}

	// Should return the same instance for method chaining
	if result != err {
		t.Error("WithType should return the same instance")
	}
}

// Test WithCommonType method
func TestRFC9457Error_WithCommonType(t *testing.T) {
	err := NewRFC9457Error(400, "Bad Request", "Invalid input")

	result := err.WithCommonType(CommonProblemTypes.ValidationError)

	if result.Type != CommonProblemTypes.ValidationError {
		t.Errorf("Type = %s, want %s", result.Type, CommonProblemTypes.ValidationError)
	}
}

// Test WithInstance method
func TestRFC9457Error_WithInstance(t *testing.T) {
	err := NewRFC9457Error(404, "Not Found", "User not found")
	instance := "/api/users/123"

	result := err.WithInstance(instance)

	if result.Instance != instance {
		t.Errorf("Instance = %s, want %s", result.Instance, instance)
	}

	// Should return the same instance for method chaining
	if result != err {
		t.Error("WithInstance should return the same instance")
	}
}

// Test WithExtension method
func TestRFC9457Error_WithExtension(t *testing.T) {
	err := NewRFC9457Error(400, "Bad Request", "Invalid input")

	result := err.WithExtension("user_id", "123").WithExtension("retry_after", 60)

	if result.Extensions["user_id"] != "123" {
		t.Errorf("Extensions[user_id] = %v, want 123", result.Extensions["user_id"])
	}

	if result.Extensions["retry_after"] != 60 {
		t.Errorf("Extensions[retry_after] = %v, want 60", result.Extensions["retry_after"])
	}
}

// Test WithMultipleProblems method
func TestRFC9457Error_WithMultipleProblems(t *testing.T) {
	err := NewRFC9457Error(400, "Multiple Validation Errors", "Several fields are invalid")

	problems := []interface{}{
		map[string]interface{}{"field": "email", "message": "Invalid format"},
		map[string]interface{}{"field": "age", "message": "Must be positive"},
	}

	result := err.WithMultipleProblems(problems)

	if result.Extensions["problems"] == nil {
		t.Error("Extensions[problems] should not be nil")
	}
}

// Test WithTraceID method
func TestRFC9457Error_WithTraceID(t *testing.T) {
	err := NewRFC9457Error(500, "Internal Server Error", "Database connection failed")
	traceID := "abc123def456"

	result := err.WithTraceID(traceID)

	if result.Extensions["trace-id"] != traceID {
		t.Errorf("Extensions[trace-id] = %v, want %s", result.Extensions["trace-id"], traceID)
	}
}

// Test WithRetryAfter method
func TestRFC9457Error_WithRetryAfter(t *testing.T) {
	err := NewRFC9457Error(503, "Service Unavailable", "Server is temporarily overloaded")
	retryAfter := 120

	result := err.WithRetryAfter(retryAfter)

	if result.Extensions["retry-after"] != retryAfter {
		t.Errorf("Extensions[retry-after] = %v, want %d", result.Extensions["retry-after"], retryAfter)
	}
}

// Test StatusCode method
func TestRFC9457Error_StatusCode(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"400 Bad Request", http.StatusBadRequest},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC9457Error(tt.status, "Test", "Test detail")
			if got := err.StatusCode(); got != tt.status {
				t.Errorf("StatusCode() = %v, want %v", got, tt.status)
			}
		})
	}
}

// Test ErrorMessage method
func TestRFC9457Error_ErrorMessage(t *testing.T) {
	tests := []struct {
		name   string
		detail string
	}{
		{"Simple message", "This is a test error"},
		{"Empty message", ""},
		{"Message with special characters", "Error: 특수문자 테스트 & symbols!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC9457Error(400, "Test", tt.detail)
			if got := err.ErrorMessage(); got != tt.detail {
				t.Errorf("ErrorMessage() = %v, want %v", got, tt.detail)
			}
		})
	}
}

// Test IsCommonType method
func TestRFC9457Error_IsCommonType(t *testing.T) {
	tests := []struct {
		name     string
		typeURI  string
		expected bool
	}{
		{"Validation Error", CommonProblemTypes.ValidationError, true},
		{"Authentication Required", CommonProblemTypes.AuthenticationRequired, true},
		{"Custom Type", "https://example.com/custom-error", false},
		{"About Blank", "about:blank", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC9457Error(400, "Test", "Test detail").WithType(tt.typeURI)
			if got := err.IsCommonType(); got != tt.expected {
				t.Errorf("IsCommonType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test IsDereferenceable method
func TestRFC9457Error_IsDereferenceable(t *testing.T) {
	tests := []struct {
		name     string
		typeURI  string
		expected bool
	}{
		{"HTTP URL", "http://example.com/problem", true},
		{"HTTPS URL", "https://example.com/problem", true},
		{"URN", "urn:problem:validation-error", false},
		{"About Blank", "about:blank", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRFC9457Error(400, "Test", "Test detail").WithType(tt.typeURI)
			if got := err.IsDereferenceable(); got != tt.expected {
				t.Errorf("IsDereferenceable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test MarshalJSON method
func TestRFC9457Error_MarshalJSON(t *testing.T) {
	err := NewRFC9457Error(400, "Bad Request", "Invalid input").
		WithType("https://example.com/problems/validation").
		WithInstance("/api/users/123").
		WithExtension("user_id", "123").
		WithExtension("errors", []string{"email is required", "age must be positive"})

	jsonBytes, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal RFC9457Error: %v", marshalErr)
	}

	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(jsonBytes, &result); unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", unmarshalErr)
	}

	// Check standard fields
	if result["type"] != "https://example.com/problems/validation" {
		t.Errorf("type = %v, want https://example.com/problems/validation", result["type"])
	}

	if result["title"] != "Bad Request" {
		t.Errorf("title = %v, want Bad Request", result["title"])
	}

	if result["status"] != float64(400) {
		t.Errorf("status = %v, want 400", result["status"])
	}

	if result["detail"] != "Invalid input" {
		t.Errorf("detail = %v, want Invalid input", result["detail"])
	}

	if result["instance"] != "/api/users/123" {
		t.Errorf("instance = %v, want /api/users/123", result["instance"])
	}

	// Check extensions
	if result["user_id"] != "123" {
		t.Errorf("user_id = %v, want 123", result["user_id"])
	}

	if result["errors"] == nil {
		t.Error("errors extension should be present")
	}
}

// Test Validate method
func TestRFC9457Error_Validate(t *testing.T) {
	tests := []struct {
		name    string
		err     *RFC9457Error
		wantErr bool
	}{
		{
			name:    "Valid error",
			err:     NewRFC9457Error(400, "Bad Request", "Invalid input"),
			wantErr: false,
		},
		{
			name: "Missing status",
			err: &RFC9457Error{
				Title:  "Bad Request",
				Detail: "Invalid input",
			},
			wantErr: true,
		},
		{
			name:    "Invalid status code (too low)",
			err:     NewRFC9457Error(99, "Invalid", "Invalid status"),
			wantErr: true,
		},
		{
			name:    "Invalid status code (too high)",
			err:     NewRFC9457Error(600, "Invalid", "Invalid status"),
			wantErr: true,
		},
		{
			name:    "Invalid type URI",
			err:     NewRFC9457Error(400, "Bad Request", "Invalid input").WithType("invalid-uri"),
			wantErr: true,
		},
		{
			name:    "Valid about:blank type",
			err:     NewRFC9457Error(400, "Bad Request", "Invalid input").WithType("about:blank"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.err.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test ToHttpError method
func TestRFC9457Error_ToHttpError(t *testing.T) {
	// Create a problem detail
	problem := NewRFC9457Error(403, "Forbidden", "Access denied")
	problem = problem.WithType("https://example.com/errors/forbidden")
	problem = problem.WithInstance("/api/resource/123")
	problem = problem.WithExtension("user_id", "abc123")

	// Convert to HttpError
	httpErr := problem.ToHttpError()

	// Check status code
	if httpErr.StatusCode() != 403 {
		t.Errorf("StatusCode() = %d, want 403", httpErr.StatusCode())
	}

	// Check that the message contains JSON
	message := httpErr.ErrorMessage()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(message), &result); err != nil {
		t.Errorf("Error message should be valid JSON: %v", err)
	}

	// Verify JSON content
	if result["type"] != "https://example.com/errors/forbidden" {
		t.Errorf("JSON type = %v, want https://example.com/errors/forbidden", result["type"])
	}
}

// Test ToRFC7807Error method
func TestRFC9457Error_ToRFC7807Error(t *testing.T) {
	// Create RFC9457 error with extensions
	rfc9457 := NewRFC9457Error(400, "Bad Request", "Invalid input").
		WithType("https://example.com/problems/validation").
		WithInstance("/api/users/123").
		WithExtension("user_id", "123").
		WithExtension("retry_after", 60)

	// Convert to RFC7807
	rfc7807 := rfc9457.ToRFC7807Error()

	// Check basic fields
	if rfc7807.Type != rfc9457.Type {
		t.Errorf("Type = %s, want %s", rfc7807.Type, rfc9457.Type)
	}

	if rfc7807.Title != rfc9457.Title {
		t.Errorf("Title = %s, want %s", rfc7807.Title, rfc9457.Title)
	}

	if rfc7807.Status != rfc9457.Status {
		t.Errorf("Status = %d, want %d", rfc7807.Status, rfc9457.Status)
	}

	if rfc7807.Detail != rfc9457.Detail {
		t.Errorf("Detail = %s, want %s", rfc7807.Detail, rfc9457.Detail)
	}

	if rfc7807.Instance != rfc9457.Instance {
		t.Errorf("Instance = %s, want %s", rfc7807.Instance, rfc9457.Instance)
	}

	// Check extensions
	if rfc7807.Extensions["user_id"] != "123" {
		t.Errorf("Extensions[user_id] = %v, want 123", rfc7807.Extensions["user_id"])
	}

	if rfc7807.Extensions["retry_after"] != 60 {
		t.Errorf("Extensions[retry_after] = %v, want 60", rfc7807.Extensions["retry_after"])
	}
}

// Test Problem Helper Functions
func TestProblemHelperFunctions9457(t *testing.T) {
	tests := []struct {
		name      string
		fn        func(string, ...string) *RFC9457Error
		status    int
		detail    string
		wantTitle string
	}{
		{
			name:      "BadRequestProblem9457",
			fn:        BadRequestProblem9457,
			status:    http.StatusBadRequest,
			detail:    "Invalid input",
			wantTitle: "Bad Request",
		},
		{
			name:      "BadRequestProblem9457 with custom title",
			fn:        BadRequestProblem9457,
			status:    http.StatusBadRequest,
			detail:    "Invalid input",
			wantTitle: "Validation Error",
		},
		{
			name:      "NotFoundProblem9457",
			fn:        NotFoundProblem9457,
			status:    http.StatusNotFound,
			detail:    "Resource not found",
			wantTitle: "Not Found",
		},
		{
			name:      "NotFoundProblem9457 with custom title",
			fn:        NotFoundProblem9457,
			status:    http.StatusNotFound,
			detail:    "Resource not found",
			wantTitle: "Missing Resource",
		},
		{
			name:      "ForbiddenProblem9457",
			fn:        ForbiddenProblem9457,
			status:    http.StatusForbidden,
			detail:    "Access denied",
			wantTitle: "Forbidden",
		},
		{
			name:      "UnauthorizedProblem9457",
			fn:        UnauthorizedProblem9457,
			status:    http.StatusUnauthorized,
			detail:    "Authentication required",
			wantTitle: "Unauthorized",
		},
		{
			name:      "InternalServerErrorProblem9457",
			fn:        InternalServerErrorProblem9457,
			status:    http.StatusInternalServerError,
			detail:    "Something went wrong",
			wantTitle: "Internal Server Error",
		},
		{
			name:      "ConflictProblem9457",
			fn:        ConflictProblem9457,
			status:    http.StatusConflict,
			detail:    "Resource already exists",
			wantTitle: "Conflict",
		},
		{
			name:      "TooManyRequestsProblem9457",
			fn:        TooManyRequestsProblem9457,
			status:    http.StatusTooManyRequests,
			detail:    "Rate limit exceeded",
			wantTitle: "Too Many Requests",
		},
		{
			name:      "ServiceUnavailableProblem9457",
			fn:        ServiceUnavailableProblem9457,
			status:    http.StatusServiceUnavailable,
			detail:    "Service is down",
			wantTitle: "Service Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err *RFC9457Error
			if tt.wantTitle != "Bad Request" && tt.wantTitle != "Not Found" &&
				tt.wantTitle != "Forbidden" && tt.wantTitle != "Unauthorized" &&
				tt.wantTitle != "Internal Server Error" && tt.wantTitle != "Conflict" &&
				tt.wantTitle != "Too Many Requests" && tt.wantTitle != "Service Unavailable" {
				// Custom title case
				err = tt.fn(tt.detail, tt.wantTitle)
			} else {
				// Default title case
				err = tt.fn(tt.detail)
			}

			if err.Status != tt.status {
				t.Errorf("Status = %d, want %d", err.Status, tt.status)
			}

			if err.Detail != tt.detail {
				t.Errorf("Detail = %s, want %s", err.Detail, tt.detail)
			}

			if err.Title != tt.wantTitle {
				t.Errorf("Title = %s, want %s", err.Title, tt.wantTitle)
			}

			// Verify that the error implements the error interface
			if err.Error() == "" {
				t.Error("Error() should return a non-empty string")
			}
		})
	}
}

// Test RFC9457 Complete Flow
func TestRFC9457CompleteFlow(t *testing.T) {
	// Create a comprehensive RFC9457 error with all features
	problem := BadRequestProblem9457("The request contains invalid data").
		WithInstance("/api/users/register").
		WithExtension("validation_errors", []map[string]string{
			{"field": "email", "message": "Invalid email format"},
			{"field": "age", "message": "Age must be between 18 and 100"},
		}).
		WithExtension("request_id", "req_123456789").
		WithTraceID("trace_abc123def456").
		WithRetryAfter(0) // No retry for validation errors

	// Test error interface
	errorMsg := problem.Error()
	expectedErrorMsg := "400: Bad Request - The request contains invalid data"
	if errorMsg != expectedErrorMsg {
		t.Errorf("Error() = %s, want %s", errorMsg, expectedErrorMsg)
	}

	// Test status code
	if problem.StatusCode() != 400 {
		t.Errorf("StatusCode() = %d, want 400", problem.StatusCode())
	}

	// Test error message
	if problem.ErrorMessage() != "The request contains invalid data" {
		t.Errorf("ErrorMessage() = %s, want 'The request contains invalid data'", problem.ErrorMessage())
	}

	// Test JSON marshaling
	jsonBytes, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("Failed to marshal problem: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify all fields are present
	expectedFields := []string{"type", "title", "status", "detail", "instance",
		"validation_errors", "request_id", "trace-id", "retry-after"}

	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Field %s is missing from JSON output", field)
		}
	}

	// Test validation
	if err := problem.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Test conversion to HttpError
	httpErr := problem.ToHttpError()
	if httpErr.StatusCode() != 400 {
		t.Errorf("HttpError status = %d, want 400", httpErr.StatusCode())
	}

	// Test conversion to RFC7807
	rfc7807 := problem.ToRFC7807Error()
	if rfc7807.Status != 400 {
		t.Errorf("RFC7807Error status = %d, want 400", rfc7807.Status)
	}

	// Test type checking methods
	if !problem.IsCommonType() {
		t.Error("Should be recognized as common type")
	}

	if !problem.IsDereferenceable() {
		t.Error("Should be recognized as dereferenceable")
	}

	t.Log("RFC9457 complete flow test passed successfully")
}
